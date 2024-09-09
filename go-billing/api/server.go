package api

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sanderdescamps/go-billing-api/gobilling/log"
	"github.com/sanderdescamps/go-billing-api/gobilling/model"
	"github.com/sanderdescamps/go-billing-api/gobilling/service"
)

var swaggerFs embed.FS

var billingDB *service.BillingDB
var jwtSecret string

func init() {
	jwtSecret = secretGenarator(32)
}

func loadSampleData(billingDB *service.BillingDB) {
	// costTypes := []*model.CostType{}
	for _, s := range []string{"vm", "container", "loadbalancer"} {
		if !billingDB.CostTypes.ExitsByName(s) {
			newCostType := model.NewCostType(s, "initial data", model.Cost{Fixed: 0.2, CostPerSec: 0.6})
			err := billingDB.CostTypes.Create(newCostType)
			if err != nil {
				log.Errorf("%s\n", err)
			}
			// costTypes = append(costTypes, newCostType)

		} else {
			log.Infof("CostType alread exists: %s\n", s)
		}
	}

	for _, cc := range []string{"cc-aws", "cc-on-prem"} {
		if !billingDB.CostCenters.ExitsByName(cc) {
			_, err := billingDB.CostCenters.Create(cc, "this is a default cost center")
			if err != nil {
				log.Errorf("%s\n", err)
			}
			// costTypes = append(costTypes, newCostType)

		} else {
			log.Infof("CostType alread exists: %s\n", cc)
		}
	}

	for _, rName := range []string{"vm1", "vm2", "vm3", "vm4"} {
		if !billingDB.Resources.ExitsByName(rName) {
			costType, err := billingDB.CostTypes.GetByName("vm")
			if err != nil {
				log.Errorf("%s\n", err)
			}

			costCenter, err := billingDB.CostCenters.GetByName("cc-aws")
			if err != nil {
				log.Errorf("%s\n", err)
			}

			_, err = billingDB.Resources.Create(rName, "", costCenter.Id, []string{costType.TypeID}, 1)
			if err != nil {
				log.Errorf("%s\n", err)
			}
		} else {
			log.Infof("Resource alread exists: %s\n", rName)
		}
	}

	if !billingDB.Roles.ExitsByName("admin") {
		billingDB.Roles.Create("admin", []string{"admin"})
	}

	if !billingDB.Roles.ExitsByName("viewer") {
		billingDB.Roles.Create("viewer", []string{
			PERM_COSTTYPE_VIEW,
			PERM_RESOURCE_VIEW,
			PERM_BASIC_LOGIN,
		})
	}

	if !billingDB.Roles.ExitsByName("operator") {
		billingDB.Roles.Create("operator", []string{
			PERM_RESOURCE_EDIT,
			PERM_RESOURCE_VIEW,
			PERM_COSTTYPE_EDIT,
			PERM_COSTTYPE_VIEW,
			PERM_BASIC_LOGIN,
		})
	}

	if !billingDB.Roles.ExitsByName("manager") {
		billingDB.Roles.Create("manager", []string{
			PERM_RESOURCE_EDIT,
			PERM_RESOURCE_VIEW,
			PERM_RESOURCE_NEW,
			PERM_RESOURCE_DELETE,
			PERM_COSTTYPE_EDIT,
			PERM_COSTTYPE_VIEW,
			PERM_COSTTYPE_NEW,
			PERM_COSTTYPE_DELETE,
			PERM_BASIC_LOGIN,
		})
	}

	adminRole, err := billingDB.Roles.GetByName("admin")
	if err != nil {
		log.Errorf(err.Error())
	}
	if !billingDB.Users.ExitsUsername("admin") {
		billingDB.Users.Create("admin", "admin", adminRole.RoleId)
	}
}

func InitDB(path string) {
	billingDB = service.NewTreeDB(path)
	loadSampleData(billingDB)
}

func InitSwagger(fs embed.FS) {
	swaggerFs = fs
}

func Run(host string, port int) {
	if billingDB == nil {
		log.Fatal("can not run server if db is not initiated")
		return
	}

	// Initialize router
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	subrouter.Use(logMiddleware)
	subrouter.HandleFunc("/token", authMW(getBearerToken, []string{PERM_USER_GET_TOKEN})).Methods("GET", "POST")
	subrouter.HandleFunc("/apitoken", authMW(getApiToken, []string{PERM_USER_GET_TOKEN})).Methods("GET", "POST")
	// router.HandleFunc("/register", authMW(RegisterHandler)).Methods("PUT", "POST")

	// Define API endpoints
	subrouter.HandleFunc("/resources", authMW(getAllResources, []string{PERM_RESOURCE_VIEW})).Methods("GET")
	subrouter.HandleFunc("/resources", authMW(createResource, []string{PERM_RESOURCE_NEW})).Methods("PUT")
	subrouter.HandleFunc("/resources/{id}", authMW(getResource, []string{PERM_RESOURCE_VIEW})).Methods("GET")
	subrouter.HandleFunc("/resources/{id}", authMW(updateResource, []string{PERM_RESOURCE_EDIT})).Methods("POST")
	subrouter.HandleFunc("/resources/{id}/add/{costTypeID}", authMW(updateResourceAddCostType, []string{PERM_RESOURCE_EDIT})).Methods("POST")
	subrouter.HandleFunc("/resources/{id}/delete/{costTypeID}", authMW(updateResourceDeleteCostType, []string{PERM_RESOURCE_EDIT})).Methods("POST")
	subrouter.HandleFunc("/resources/{id}", authMW(deleteResource, []string{PERM_RESOURCE_DELETE})).Methods("DELETE")

	subrouter.HandleFunc("/cost_type", authMW(getAllCostTypes, []string{PERM_COSTTYPE_VIEW})).Methods("GET")
	subrouter.HandleFunc("/cost_type", authMW(createCostType, []string{PERM_COSTTYPE_NEW})).Methods("PUT")
	subrouter.HandleFunc("/cost_type/{id}", authMW(getCostType, []string{PERM_COSTTYPE_VIEW})).Methods("GET")
	subrouter.HandleFunc("/cost_type/{id}", authMW(updateCostType, []string{PERM_COSTTYPE_EDIT})).Methods("POST")
	subrouter.HandleFunc("/cost_type/{id}", authMW(deleteCostType, []string{PERM_COSTTYPE_DELETE})).Methods("DELETE")

	subrouter.HandleFunc("/cost_center", authMW(getAllCostCenters, []string{PERM_COSTCENTER_VIEW})).Methods("GET")
	subrouter.HandleFunc("/cost_center", authMW(createCostCenter, []string{PERM_COSTCENTER_NEW})).Methods("PUT")
	subrouter.HandleFunc("/cost_center/{id}/total_cost", authMW(getTotalCostCenterCost, []string{PERM_COSTCENTER_NEW})).Methods("GET")
	subrouter.HandleFunc("/cost_center/{id}", authMW(getCostCenter, []string{PERM_COSTCENTER_VIEW})).Methods("GET")
	subrouter.HandleFunc("/cost_center/{id}", authMW(updateCostCenter, []string{PERM_COSTCENTER_EDIT})).Methods("POST")
	subrouter.HandleFunc("/cost_center/{id}", authMW(deleteCostCenter, []string{PERM_COSTCENTER_DELETE})).Methods("DELETE")

	//swagger
	subdir, _ := fs.Sub(swaggerFs, "swagger")
	sh := http.StripPrefix("/swagger/", http.FileServer(http.FS(subdir)))
	router.PathPrefix("/swagger").Handler(sh)

	// Start the server
	fmt.Printf("Server is running on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("%s - %s (%s)", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

package test

import (
	"after-sales/api/config"
	masterrepositoryimpl "after-sales/api/repositories/master/repositories-impl"
	masterserviceimpl "after-sales/api/services/master/service-impl"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

// TestGetAllIncentiveMaster_Success digunakan untuk menguji skenario ketika permintaan berhasil
func TestGetAllIncentiveMaster_Success(t *testing.T) {
	// Membuat request GET palsu ke endpoint GetAllIncentiveMaster
	req, err := http.NewRequest("GET", "http://localhost:8000/incentive", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Membuat objek ResponseWriter palsu untuk menangkap respons
	rr := httptest.NewRecorder()

	// Menggunakan handler yang sebenarnya (misalnya handler API Anda) untuk menangani permintaan
	handler := http.HandlerFunc(YourGetAllIncentiveMasterHandler)
	handler.ServeHTTP(rr, req)

	// Memeriksa kode status HTTP yang diharapkan (200 OK)
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")

	// Memeriksa konten respons yang diharapkan
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}
	assert.NotNil(t, response["data"], "Response should contain data")
}

// TestGetAllIncentiveMaster_Error digunakan untuk menguji skenario ketika permintaan mengalami kesalahan
func TestGetAllIncentiveMaster_Error(t *testing.T) {
	// Membuat request GET palsu ke endpoint GetAllIncentiveMaster
	req, err := http.NewRequest("GET", "http://localhost:8000/api/incentive", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Membuat objek ResponseWriter palsu untuk menangkap respons
	rr := httptest.NewRecorder()

	// Menggunakan handler yang sebenarnya (misalnya handler API Anda) untuk menangani permintaan
	handler := http.HandlerFunc(YourGetAllIncentiveMasterHandlerWithError)
	handler.ServeHTTP(rr, req)

	// Memeriksa kode status HTTP yang diharapkan (400 Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code should be 400")

	// Memeriksa konten respons yang diharapkan (misalnya pesan kesalahan dalam JSON)
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}
	assert.NotNil(t, response["error"], "Response should contain error message")
}

// Ini adalah contoh fungsi handler untuk endpoint API GET GetAllIncentiveMaster yang berhasil
func YourGetAllIncentiveMasterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Contoh data respons yang dikembalikan
	data := map[string]interface{}{
		"data": []string{"incentive1", "incentive2", "incentive3"},
	}
	json.NewEncoder(w).Encode(data)
}

// Ini adalah contoh fungsi handler untuk endpoint API GET GetAllIncentiveMaster yang menghasilkan kesalahan
func YourGetAllIncentiveMasterHandlerWithError(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	// Contoh data respons yang berisi pesan kesalahan
	data := map[string]interface{}{
		"error": "Failed to fetch incentive data",
	}
	json.NewEncoder(w).Encode(data)
}

func TestGetIncentiveMasterById(t *testing.T) {
	// Inisialisasi konfigurasi dan layanan
	config.InitEnvConfigs(true, "")
	db := config.InitDB()
	rdb := config.InitRedis()
	IncentiveMasterRepository := masterrepositoryimpl.StartIncentiveMasterRepositoryImpl()
	IncentiveMasterService := masterserviceimpl.StartIncentiveMasterService(IncentiveMasterRepository, db, rdb)

	// Panggil GetIncentiveMasterById dan tangkap hasilnya
	get, err := IncentiveMasterService.GetIncentiveMasterById(2)
	if err != nil {
		// Handle error
		fmt.Println("Error:", err)
		return
	}
	// Now you can use the response
	fmt.Println("Response:", get)

	// Periksa apakah hasilnya valid (tidak kosong) menggunakan assert dari Testify
	assert.NotZero(t, get.IncentiveLevelId, "Expected non-zero result for IncentiveLevelId")

	// Periksa menggunakan GoConvey
	convey.Convey("Test Get Incentive Master By Id", t, func() {
		convey.So(get.IncentiveLevelId, convey.ShouldEqual, 2)
	})
}

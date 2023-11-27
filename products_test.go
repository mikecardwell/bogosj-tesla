package tesla

import (
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

var (
	ProductsJSON = `{"response":[{"id":1234567890123456,"user_id":999999,"vehicle_id":999999999999,"vin":"ABCDEFGH9AB999999","color":null,"access_type":"OWNER","display_name":"foo","option_codes":null,"cached_data":"<redacted>","granular_access":{"hide_private":false},"tokens":["token1","token2"],"state":"online","in_service":false,"id_s":"9999999999999999","calendar_enabled":true,"api_version":67,"backseat_token":null,"backseat_token_updated_at":null,"ble_autopair_enrolled":false,"vehicle_config":{"aux_park_lamps":"Eu","badge_version":0,"can_accept_navigation_requests":true,"can_actuate_trunks":true,"car_special_type":"base","car_type":"modely","charge_port_type":"CCS","cop_user_set_temp_supported":false,"dashcam_clip_save_supported":true,"default_charge_to_max":false,"driver_assist":"TeslaAP3","ece_restrictions":true,"efficiency_package":"MY2020","eu_vehicle":true,"exterior_color":"PearlWhite","exterior_trim":"Black","exterior_trim_override":"","has_air_suspension":false,"has_ludicrous_mode":false,"has_seat_cooling":false,"headlamp_type":"Global","interior_trim_type":"Black2","key_version":2,"motorized_charge_port":true,"paint_color_override":"20,20,20,0.01,0.04","performance_package":"Base","plg":true,"pws":true,"rear_drive_unit":"PM216MOSFET","rear_seat_heaters":1,"rear_seat_type":0,"rhd":false,"roof_color":"RoofColorGlass","seat_type":null,"spoiler_type":"None","sun_roof_installed":null,"supports_qr_pairing":false,"third_row_seats":"None","timestamp":1700733459821,"trim_badging":"74d","use_range_badging":true,"utc_offset":3600,"webcam_selfie_supported":true,"webcam_supported":true,"wheel_type":"Apollo19"},"command_signing":"required","release_notes_supported":true},{"energy_site_id":12345678901234,"resource_type":"battery","site_name":"my site","id":"STE19700101-00001","gateway_id":"9999999-01-D--TG11111111111F","asset_site_id":"redacted-uuid","warp_site_number":"STE19700101-00001","energy_left":0,"total_pack_energy":12981,"percentage_charged":0,"battery_type":"ac_powerwall","battery_power":-20,"go_off_grid_test_banner_enabled":null,"storm_mode_enabled":true,"powerwall_onboarding_settings_set":true,"powerwall_tesla_electric_interested_in":null,"vpp_tour_enabled":null,"sync_grid_alert_enabled":false,"breaker_alert_enabled":true,"components":{"battery":true,"battery_type":"ac_powerwall","solar":true,"solar_type":"pv_panel","grid":true,"load_meter":true,"market_type":"residential"},"features":{"rate_plan_manager_no_pricing_constraint":true}}],"count":2}`
)

func TestProductsSpec(t *testing.T) {
	ts := serveHTTP(t)
	defer ts.Close()

	client := NewTestClient(ts)

	c.Convey("Should get products", t, func() {
		products, err := client.Products()
		c.So(err, c.ShouldBeNil)
		c.So((string)(products[0].ID), c.ShouldEqual, "1234567890123456")
		c.So((string)(products[1].ID), c.ShouldEqual, "STE19700101-00001")
	})
}

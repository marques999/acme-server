package admin

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/auth"
	"github.com/marques999/acme-server/customers"
	"github.com/marques999/acme-server/orders"
	"github.com/marques999/acme-server/products"
	"github.com/pborman/uuid"
)

func clearTables(database *gorm.DB) error {

	if ex := database.Delete(orders.Order{}).Error; ex != nil {
		return ex
	} else if ex := database.Delete(customers.Customer{}).Error; ex != nil {
		return ex
	} else if ex := database.Delete(products.Product{}).Error; ex != nil {
		return ex
	} else {
		return database.Delete(customers.CreditCard{}).Error
	}
}

func populateTables(database *gorm.DB) error {

	user1Password, _ := auth.GeneratePassword("admin")
	user2Password, _ := auth.GeneratePassword("acmetest")
	user3Password, _ := auth.GeneratePassword("acmetest")

	user1 := customers.Customer{
		Address:   "R. Dr. Roberto Frias 291",
		Country:   "PT",
		Email:     "admin@acme.pt",
		Name:      "Administrator",
		Password:  string(user1Password),
		TaxNumber: "123456789",
		Username:  "admin",
		CreditCard: &customers.CreditCard{
			Number:   "123456789",
			Type:     "VISA",
			Validity: time.Now().AddDate(5, 0, 0),
		},
	}

	if ex := database.Save(&user1).Error; ex != nil {
		return ex
	}

	user2 := customers.Customer{
		Address:   "Rua Costa, nº 176",
		Country:   "PT",
		Email:     "up201305642@fe.up.pt",
		Name:      "Diogo Marques",
		Password:  string(user2Password),
		TaxNumber: "222555777",
		Username:  "marques999",
		CreditCard: &customers.CreditCard{
			Number:   "800200400",
			Validity: time.Now().AddDate(3, 6, 0),
			Type:     "Maestro",
		},
	}

	if ex := database.Save(&user2).Error; ex != nil {
		return ex
	}

	user3 := customers.Customer{
		Address:   "",
		Country:   "PT",
		Email:     "up201303930@fe.up.pt",
		Name:      "José Teixeira",
		Password:  string(user3Password),
		TaxNumber: "987654321",
		Username:  "jabst",
		CreditCard: &customers.CreditCard{
			Number:   "360420999",
			Validity: time.Now().AddDate(1, 3, 13),
			Type:     "Memez",
		},
	}

	if ex := database.Save(&user3).Error; ex != nil {
		return ex
	}

	product1 := products.Product{
		Brand:   "Acer",
		Name:    "Aspire E5-571G-72M5",
		Barcode: "5701928664431",
		Price:   490.00,
		Description: "Aspire E Series laptops are great choices for everyday" +
			"users, with lots of appealing options and an attractive design" +
			"that exceed expectations. With many enhanced components, color " +
			"choices, and a textured metallic finish, the Aspire E makes " +
			"everyday better.",
	}

	if ex := database.Save(&product1).Error; ex != nil {
		return ex
	}

	product2 := products.Product{
		Brand:   "Cooler Master",
		Name:    "MasterKeys Lite L Combo",
		Barcode: "4719512052914",
		Price:   54.99,
		Description: "Mem-chanical Switches: Cooler Master’s exclusive switches " +
			"are durable and feel like mechanical switches with satisfying" +
			"tactile feedback. Brillant Illumination - Zoned RGB backlighting " +
			"system with multiple lighting effects. 26-Key Anti-Ghosting - " +
			"Ensure each key press is correctly detected regardless how fast " +
			"and furious it gets. Cherry MX Compatible - Customize each key and " +
			"fully express yourself.Precision Optical Sensor - Your cursor goes " +
			"where you want, as fast as you want, even during intense gaming. " +
			"Be in total control.",
	}

	if ex := database.Save(&product2).Error; ex != nil {
		return ex
	}

	product3 := products.Product{
		Brand:   "MSI",
		Name:    "GeForce GTX 1060 Gaming X 6GB",
		Barcode: "4719072470364",
		Price:   339.00,
		Description: "*THE ULTIMATE PC GAMING PLATFORM*\nGeForce GTX graphics cards are the most " +
			"advanced ever created. Discover unprecedented performance, power efficiency, and " +
			"next-generation gaming experiences.\n*Nvidia VR READY*\nDiscover next-generation VR " +
			"performance, the lowest latency, and plug-and-play compatibility with leading headsets " +
			"driven by NVIDIA VRWorks™ technologies. VR audio, physics, and haptics let you hear " +
			"and feel every moment.\n*THE LATEST GAMING TECHNOLOGIES*\nPascal is built to meet the " +
			"demands of next generation displays, including VR, ultra-high-resolution, and multiple " +
			"monitors. It features NVIDIA GameWorks™ technologies for extremely smooth gameplay and " +
			"cinematic experiences.Plus, it includes revolutionary new 360-degree image capture.\n" +
			"*PERFORMANCE*\nPascal - powered graphics cards give you superior performance and power " +
			"efficiency, built using ultra-fast FinFET and supporting DirectX™ 12 features to deliver " +
			"the fastest, smoothest, most power-efficient gaming experiences.",
	}

	if ex := database.Save(&product3).Error; ex != nil {
		return ex
	}

	product4 := products.Product{
		Brand:   "Intel",
		Name:    "i5-6600K 3.5GHz 6MB Sk1151",
		Barcode: "5032037076142",
		Price:   279.90,
		Description: "The Intel Core i5-6600K is based on the new \"Skylake\" 14nm manufacturing " +
			"process. Sporting 4 physical cores with base/turbo clocks of 3.5/3.9 GHz the 6600K and " +
			"its predecessor, the 4690K share the same basic configuration and disappointingly, " +
			"offer similar performance.",
	}

	if ex := database.Save(&product4).Error; ex != nil {
		return ex
	}

	product5 := products.Product{
		Brand:   "Asus",
		Name:    "Z170 Pro Gaming",
		Barcode: "",
		Price:   164.91,
		Description: "High-value, feature-packed, performance-optimized Z170 ATX board LGA1151 " +
			"socket for 6th Gen Intel® Core™ Desktop Processors.\n- Dual DDR4 3400 (OC) support\n- " +
			"PRO Clock technology, 5-Way Optimization and 2nd-generation T-Topology: Easy and stable " +
			"overclocking\n- SupremeFX: Flawless audio that makes you part of the game\n- Intel " +
			"Gigabit Ethernet, LANGuard & GameFirst III: Top-speed protected networking\n- RAMCache: " +
			"Speed up your game loads\n- USB 3.1 Type A/C & M.2: Ultra-speedy transfers for faster " +
			"gaming\n- Gamer's Guardian: Highly-durable components and smart DIY features\n- Sonic " +
			"Radar ll: Scan and detect your enemies to dominate",
	}

	if ex := database.Save(&product5).Error; ex != nil {
		return ex
	}

	if ex := database.Save(&orders.Order{
		Customer: &user2,
		Valid:    true,
		Token:    uuid.NewUUID(),
		Products: []products.Product{product1},
	}).Error; ex != nil {
		return ex
	}

	if ex := database.Save(&orders.Order{
		Valid:    false,
		Customer: &user1,
		Products: []products.Product{product2},
	}).Error; ex != nil {
		return ex
	}

	if ex := database.Save(&orders.Order{
		Customer: &user2,
		Valid:    true,
		Token:    uuid.NewUUID(),
		Products: []products.Product{product4, product5},
	}).Error; ex != nil {
		return ex
	}

	return database.Save(&orders.Order{
		Customer: &user3,
		Valid:    true,
		Token:    uuid.NewUUID(),
		Products: []products.Product{product3},
	}).Error
}

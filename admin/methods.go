package admin

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/auth"
	"github.com/marques999/acme-server/customers"
	"github.com/marques999/acme-server/orders"
	"github.com/marques999/acme-server/products"
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
		Username:  "admin",
		Name:      "Administrator",
		Password:  string(user1Password),
		TaxNumber: "123456789",
		PublicKey: `-----BEGIN PUBLIC KEY-----
MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvAL1L9h1N9xqNe0I4ddyjKD6lv0ArcEhBJbU550urvmvJ
qa1Rm8Zr+V0+VCp9swcCAwEAAQ==
-----END PUBLIC KEY-----`,
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
		Username:  "marques999",
		Name:      "Diogo Marques",
		Password:  string(user2Password),
		TaxNumber: "222555777",
		PublicKey: `-----BEGIN PUBLIC KEY-----
MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvAL1L9h1N9xqNe0I4ddyjKD6lv0ArcEhBJbU550urvmvJ
qa1Rm8Zr+V0+VCp9swcCAwEAAQ==
-----END PUBLIC KEY-----`,
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
		Username:  "jabst",
		Name:      "José Teixeira",
		Password:  string(user3Password),
		TaxNumber: "987654321",
		PublicKey: `-----BEGIN PUBLIC KEY-----
MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvAL1L9h1N9xqNe0I4ddyjKD6lv0ArcEhBJbU550urvmvJ
qa1Rm8Zr+V0+VCp9swcCAwEAAQ==
-----END PUBLIC KEY-----`,
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
		Brand:    "Acer",
		Name:     "Aspire E5-571G-72M5",
		Price:    490.00,
		Barcode:  "4713147489589",
		ImageUri: "https://www.notebookcheck.net/fileadmin/Notebooks/Acer/Aspire_E5-571G-536E/Aspire_E5_571_531_551_521_511_nontouch_black_glare_gallery_01.png",
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
		Brand:    "Cooler Master",
		Name:     "MasterKeys Lite L Combo",
		Price:    54.99,
		Barcode:  "4719512052914",
		ImageUri: "http://cdn1.bigcommerce.com/server3900/9dd4a/products/1287/images/6402/Cooler_Master_MasterKeys_Lite_L_Combo_10__09359.1468418659.1280.1280.jpg",
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
		Brand:    "MSI",
		Name:     "GeForce GTX 1060 Gaming X 6GB",
		Price:    339.00,
		Barcode:  "4719072470364",
		ImageUri: "https://www.picclickimg.com/00/s/OTU4WDEyODA=/z/pBMAAOSwvflZO0rM/$/MSI-GTX-1060-GAMING-X-6G-Nvidia-GeForce-_1.jpg",
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
		Brand:    "Intel",
		Name:     "i5-6600K 3.5GHz 6MB Sk1151",
		Price:    279.90,
		Barcode:  "5032037076142",
		ImageUri: "https://images10.newegg.com/ProductImage/19-117-561-02.jpg",
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
		Barcode: "4712900114874",
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

	product6 := products.Product{
		Brand:    "Motorola",
		Name:     "Moto G5 Plus 5.2\" 32GB Dual SIM",
		Price:    279.90,
		Barcode:  "6947681540651",
		ImageUri: "http://www.mobilewithprices.com/products/motorola-moto-g5-Plus.jpg",
		Description: "* OUTSTANDING CAMERAS *\nThe 12 MP rear camera focuses up to 60% faster than " +
			"ever before. Switch to the wide-angle front camera for group selfies.\n* PRECISION-CRAFTED " +
			"METAL DESIGN *\nOne of the first new Moto G phones made from high-grade aluminum, it " +
			"looks as great as it performs.\n* FUEL UP FAST *\nAll-day battery and up to 6 hours of " +
			"battery life in just 15 minutes with TurboPower charging.\n* FAST OCTA-CORE PROCESSOR\n" +
			"Apps run smoothly thanks to a blazing-fast Qualcomm® Snapdragon™ 2.0 GHz octa-core " +
			"processor.\n* FINGERPRINT READER *\nInstantly unlock your phone. No passcode required.\n" +
			"* MOTO EXPERIENCES *\nGet shortcuts to the features you use most, like turning on the " +
			"camera with a twist of your wrist.",
	}

	if ex := database.Save(&product6).Error; ex != nil {
		return ex
	}

	order1 := &orders.Order{
		Status:   1,
		Customer: user2.ID,
		Products: []products.Product{product1},
	}

	if ex := database.Save(&order1).Error; ex != nil {
		return ex
	}

	order1.Token, _ = orders.GenerateToken(order1)

	if ex := database.Save(&order1).Error; ex != nil {
		return ex
	}

	order2 := &orders.Order{
		Status:   0,
		Customer: user1.ID,
		Products: []products.Product{product2},
	}

	if ex := database.Save(&order2).Error; ex != nil {
		return ex
	}

	order3 := &orders.Order{
		Status:   1,
		Customer: user2.ID,
		Products: []products.Product{product4, product5},
	}

	if ex := database.Save(&order3).Error; ex != nil {
		return ex
	}

	order3.Token, _ = orders.GenerateToken(order3)

	if ex := database.Save(&order3).Error; ex != nil {
		return ex
	}

	order4 := &orders.Order{
		Status:   1,
		Customer: user3.ID,
		Products: []products.Product{product3},
	}

	order4.Token, _ = orders.GenerateToken(order4)

	return database.Save(&order4).Error
}
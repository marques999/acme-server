package admin

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/marques999/acme-server/auth"
	"github.com/marques999/acme-server/customers"
	"github.com/marques999/acme-server/orders"
	"github.com/marques999/acme-server/products"
)

func generatePassword(password string) string {
	hashedPassword, _ := auth.GeneratePassword(password)
	return hashedPassword
}

func clearTables(db *gorm.DB) error {

	if ex := db.Delete(orders.Order{}).Error; ex != nil {
		return ex
	} else if ex := db.Delete(customers.Customer{}).Error; ex != nil {
		return ex
	} else if ex := db.Delete(products.Product{}).Error; ex != nil {
		return ex
	} else {
		return db.Delete(customers.CreditCard{}).Error
	}
}

func registerOrder(db *gorm.DB, customer customers.Customer, products []products.Product, status int) error {

	order := &orders.Order{
		Status:   status,
		Products: products,
		Customer: customer.ID,
		Total:    orders.CalculateTotal(products),
	}

	if dbException := db.Save(order).Error; dbException != nil || status == orders.ValidationFailed {
		return dbException
	} else {
		order.Token, _ = orders.GenerateToken(order)
	}

	return db.Save(&order).Error
}

func populateTables(db *gorm.DB) error {

	user1 := customers.Customer{
		Name:      "Administrator",
		Username:  "admin",
		Password:  generatePassword("admin"),
		TaxNumber: "930248516",
		Address1:  "Rua Branco, Nº 25",
		Address2:  "8681-962 Tomar",
		Country:   "PT",
		PublicKey: `MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvAL1L9h1N9xqNe0I4ddyjKD6lv0ArcEhBJbU550urvmvJ
qa1Rm8Zr+V0+VCp9swcCAwEAAQ==`,
		CreditCard: &customers.CreditCard{
			Type:     "VISA",
			Number:   "123456789",
			Validity: time.Now().AddDate(5, 0, 0),
		},
	}

	if dbException := db.Save(&user1).Error; dbException != nil {
		return dbException
	}

	user2 := customers.Customer{
		Name:      "Diogo Marques",
		Username:  "marques999",
		Password:  generatePassword("r0wsauce"),
		TaxNumber: "761489053",
		Address1:  "Rua São Diogo, Nº 855",
		Address2:  "6311-969 Vendas Novas",
		Country:   "PT",
		PublicKey: `MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvAKCRuhMUuFoJvDVeicvyfyQf9ADQ1qNe+dabNSpOkr76
FcVTBd+TBe2sEshVefUCAwEAAQ==`,
		CreditCard: &customers.CreditCard{
			Number:   "310867542",
			Validity: time.Now().AddDate(3, 6, 0),
			Type:     "Maestro",
		},
	}

	if dbException := db.Save(&user2).Error; dbException != nil {
		return dbException
	}

	user3 := customers.Customer{
		Username:  "jabst",
		Password:  generatePassword("bighotshaq"),
		Name:      "José Teixeira",
		TaxNumber: "685102439",
		Address1:  "Avenida Lima, Nº 167",
		Address2:  "7049-952 Santa Cruz",
		Country:   "PT",
		PublicKey: `MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvALLIEFJe1v3hiGpzYlzo/hxEXBW2XrA47b/S2i0X7ZZv
08HLhNfdPr2XC8ZzLpECAwEAAQ==`,
		CreditCard: &customers.CreditCard{
			Type:     "Mastercard",
			Number:   "360420999",
			Validity: time.Now().AddDate(1, 3, 13),
		},
	}

	if dbException := db.Save(&user3).Error; dbException != nil {
		return dbException
	}

	if dbException := db.Save(&customers.Customer{
		Username:  "somouco",
		Name:      "Carlos Samouco",
		Password:  generatePassword("skibidipap"),
		TaxNumber: "537812640",
		Address1:  "Travessa Mia Assunção, Nº 532",
		Address2:  "5334-964 Coimbra",
		Country:   "PT",
		PublicKey: `MEowDQYJKoZIhvcNAQEBBQADOQAwNgIvAK0smd9hF2yMJOeidEDq2GieQJY2Ac3bRpoXeOpiD/Oi
pBrNyqlMpzEKUF917T0CAwEAAQ==`,
		CreditCard: &customers.CreditCard{
			Type:     "VISA Electron",
			Number:   "863101278",
			Validity: time.Now().AddDate(2, 5, 5),
		},
	}).Error; dbException != nil {
		return dbException
	}

	product1 := products.Product{
		Brand:    "Acer",
		Name:     "Aspire E5-571G-72M5",
		Price:    490.00,
		Barcode:  "4713147489589",
		ImageUri: "https://www.notebookcheck.net/fileadmin/Notebooks/Acer/Aspire_E5-571G-536E/Aspire_E5_571_531_551_521_511_nontouch_black_glare_gallery_01.png",
		Description: "Aspire E Series laptops are great choices for everyday users, with lots of " +
			"appealing options and an attractive design that exceed expectations. With many enhanced " +
			"components, color choices, and a textured metallic finish, the Aspire E makes everyday better.",
	}

	if dbException := db.Save(&product1).Error; dbException != nil {
		return dbException
	}

	product2 := products.Product{
		Brand:    "Cooler Master",
		Name:     "MasterKeys Lite L Combo",
		Price:    54.99,
		Barcode:  "884102029028",
		ImageUri: "https://www.pcdiga.com/media/catalog/product/cache/1/image/2718f121925249d501c6086d4b8f9401/2/2/22674_1.jpg",
		Description: "Mem-chanical Switches: Cooler Master’s exclusive switches are durable and feel " +
			"like mechanical switches with satisfying tactile feedback. Zoned RGB backlighting system " +
			"with multiple lighting effects. 26-Key Anti-Ghosting - Ensure each key press is correctly " +
			"detected regardless how fast and furious it gets. Cherry MX Compatible - customize each " +
			"key and fully express yourself. Precision Optical Sensor - Your cursor goes where you want, " +
			"as fast as you want, even during intense gaming. Be in total control.",
	}

	if dbException := db.Save(&product2).Error; dbException != nil {
		return dbException
	}

	product3 := products.Product{
		Brand:    "MSI",
		Name:     "GeForce GTX 1060 Gaming X 6GB",
		Price:    339.00,
		Barcode:  "824142132142",
		ImageUri: "https://www.pcdiga.com/media/catalog/product/cache/1/image/2718f121925249d501c6086d4b8f9401/2/2/22264_1.jpg",
		Description: "GeForce GTX graphics cards are the most advanced ever created. Discover " +
			"unprecedented performance, power efficiency, and next-generation gaming experiences. " +
			"Discover next-generation VR performance, the lowest latency, and plug-and-play " +
			"compatibility with leading headsets driven by NVIDIA VRWorks™ technologies. VR audio, " +
			"physics, and haptics let you hear and feel every moment. Pascal is built to meet the " +
			"demands of next generation displays, including VR, ultra-high-resolution, and multiple " +
			"monitors. It features NVIDIA GameWorks™ technologies for extremely smooth gameplay and " +
			"cinematic experiences. Plus, it includes revolutionary new 360-degree image capture.\n" +
			"Pascal powered graphics cards give you superior performance and power efficiency, built " +
			"using ultra-fast FinFET and supporting DirectX™ 12 features to deliver the fastest, " +
			"smoothest, most power-efficient gaming experiences.",
	}

	if dbException := db.Save(&product3).Error; dbException != nil {
		return dbException
	}

	product4 := products.Product{
		Brand:    "Asus",
		Name:     "Z170 Pro Gaming",
		Price:    164.91,
		Barcode:  "889349114872",
		ImageUri: "https://www.asus.com/media/global/products/JIM9ojJyz4lZRy4k/P_setting_fff_1_90_end_500.png",
		Description: "High-value, feature-packed, performance-optimized Z170 ATX board LGA1151 " +
			"socket for 6th Gen Intel® Core™ Desktop Processors.\n- Dual DDR4 3400 (OC) support\n- " +
			"PRO Clock technology, 5-Way Optimization and 2nd-generation T-Topology: Easy and stable " +
			"overclocking\n- SupremeFX: Flawless audio that makes you part of the game\n- Intel " +
			"Gigabit Ethernet, LANGuard & GameFirst III: Top-speed protected networking\n- RAMCache: " +
			"Speed up your game loads\n- USB 3.1 Type A/C & M.2: Ultra-speedy transfers for faster " +
			"gaming\n- Gamer's Guardian: Highly-durable components and smart DIY features\n- Sonic " +
			"Radar ll: Scan and detect your enemies to dominate",
	}

	if dbException := db.Save(&product4).Error; dbException != nil {
		return dbException
	}

	if dbException := db.Save(&products.Product{
		Brand:    "Intel",
		Name:     "i5-6600K 3.5GHz 6MB Sk1151",
		Price:    279.90,
		Barcode:  "735858301077",
		ImageUri: "https://www.pcdiga.com/media/catalog/product/cache/1/image/2718f121925249d501c6086d4b8f9401/1/7/17863_1.jpg",
		Description: "The Intel Core i5-6600K is based on the new \"Skylake\" 14nm manufacturing " +
			"process. Sporting 4 physical cores with base/turbo clocks of 3.5/3.9 GHz the 6600K and " +
			"its predecessor, the 4690K share the same basic configuration and disappointingly, " +
			"offer similar performance.",
	}).Error; dbException != nil {
		return dbException
	}

	if dbException := db.Save(&products.Product{
		Brand:    "G.Skill",
		Name:     "Trident Z 16GB (2x8GB) DDR4-3000MHz CL15",
		Price:    192.50,
		Barcode:  "848354015451",
		ImageUri: "https://www.pcdiga.com/media/catalog/product/cache/1/image/2718f121925249d501c6086d4b8f9401/s/e/sem-t_tulo-2_8.jpg",
		Description: "Building on the strong success of G.SKILL Trident series, Trident Z series " +
			"represents one of the world’s highest performance DDR4 memory designed for the latest " +
			"6th generation Intel® Core™ processor on the Z170 series chipset. Using only the " +
			"best-in-class components and featuring dual-color construction aluminum heat-spreaders, " +
			"Trident Z series is the state-of-the-art DDR4 solution that combines performance and " +
			"beauty for PC enthusiasts and extreme overclockers to build an ultra-fast PC or achieve " +
			"new overclocking records.",
	}).Error; dbException != nil {
		return dbException
	}

	if dbException := db.Save(&products.Product{
		Brand:    "Arctic",
		Name:     "MX-4 (4g)",
		Barcode:  "872767003767",
		Price:    5.90,
		ImageUri: "https://www.pcdiga.com/media/catalog/product/cache/1/image/2718f121925249d501c6086d4b8f9401/7/1/712_1.jpg",
		Description: "ARCTIC MX-4 is a new thermal compound that guarantees exceptional heat " +
			"dissipation from components and maintains the needed stability to push your computer " +
			"system to its maximum. Similar to the performance of existing and acclaimed ARCTIC MX " +
			"series, ARCTIC MX-4 continues to be overclocker’s ultimate choice when choosing thermal " +
			"compounds. Prior to its introduction into the market, ARCTIC MX-4 has already received " +
			"a Top Product Award from PC Games Hardware in a Germany print magazine in August 2010 " +
			"as it prevailed against 11 other thermal compounds in the market. ",
	}).Error; dbException != nil {
		return dbException
	}

	if ex := registerOrder(db, user2, []products.Product{product1}, orders.ValidationComplete); ex != nil {
		return ex
	}

	if ex := registerOrder(db, user1, []products.Product{product3, product4}, orders.ValidationFailed); ex != nil {
		return ex
	}

	if ex := registerOrder(db, user2, []products.Product{product2, product4}, orders.ValidationComplete); ex != nil {
		return ex
	}

	return registerOrder(db, user3, []products.Product{product3}, orders.ValidationComplete)
}
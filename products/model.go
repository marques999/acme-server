package products

import (
	"github.com/jmoiron/sqlx"
	"github.com/marques999/acme-server/common"
)

const (
	Name        = "name"
	Brand       = "brand"
	Price       = "price"
	Barcode     = "barcode"
	Products    = "products"
	ImageUri    = "image_uri"
	Description = "description"
)

type Product struct {
	common.Model
	Name        string
	Brand       string
	Price       float64
	Barcode     string
	ImageUri    string `db:"image_uri"`
	Description string
}

type ProductJSON struct {
	Name        string  `binding:"required" json:"name"`
	Brand       string  `binding:"required" json:"brand"`
	Price       float64 `binding:"required" json:"price"`
	Barcode     string  `binding:"required" json:"barcode"`
	ImageUri    string  `binding:"required" json:"image_uri"`
	Description string  `binding:"required" json:"description"`
}

func Migrate(database *sqlx.DB) {

	if _, sqlException := database.Exec(`CREATE TABLE products(
		id serial NOT NULL CONSTRAINT products_pkey PRIMARY KEY,
		name TEXT NOT NULL,
		brand TEXT NOT NULL,
		price NUMERIC NOT NULL,
		barcode TEXT NOT NULL,
		image_uri TEXT NOT NULL,
		description TEXT NOT NULL,
		created_at timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated_at timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL)
	`); sqlException != nil {
		return
	}

	database.MustExec("CREATE INDEX IF NOT EXISTS idx_products_barcode ON products (barcode)")
	database.MustExec("CREATE UNIQUE INDEX IF NOT EXISTS uix_products_barcode ON products (barcode)")

	insertProduct(database, ProductJSON{
		Brand:    "Acer",
		Name:     "Aspire E5-571G-72M5",
		Price:    490.00,
		Barcode:  "4713147489589",
		ImageUri: "https://www.notebookcheck.net/fileadmin/Notebooks/Acer/Aspire_E5-571G-536E/Aspire_E5_571_531_551_521_511_nontouch_black_glare_gallery_01.png",
		Description: "Aspire E Series laptops are great choices for everyday users, with lots of " +
			"appealing options and an attractive design that exceed expectations. With many enhanced " +
			"components, color choices, and a textured metallic finish, the Aspire E makes everyday better.",
	})

	insertProduct(database, ProductJSON{
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
	})

	insertProduct(database, ProductJSON{
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
	})

	insertProduct(database, ProductJSON{
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
	})

	insertProduct(database, ProductJSON{
		Brand:    "Intel",
		Name:     "i5-6600K 3.5GHz 6MB Sk1151",
		Price:    279.90,
		Barcode:  "735858301077",
		ImageUri: "https://www.pcdiga.com/media/catalog/product/cache/1/image/2718f121925249d501c6086d4b8f9401/1/7/17863_1.jpg",
		Description: "The Intel Core i5-6600K is based on the new \"Skylake\" 14nm manufacturing " +
			"process. Sporting 4 physical cores with base/turbo clocks of 3.5/3.9 GHz the 6600K and " +
			"its predecessor, the 4690K share the same basic configuration and disappointingly, " +
			"offer similar performance.",
	})

	insertProduct(database, ProductJSON{
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
	})

	insertProduct(database, ProductJSON{
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
			"as it prevailed against 11 other thermal compounds in the market. ",})
}

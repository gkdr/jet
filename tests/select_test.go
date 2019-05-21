package tests

import (
	"fmt"
	. "github.com/sub0zero/go-sqlbuilder/sqlbuilder"
	"github.com/sub0zero/go-sqlbuilder/tests/.test_files/dvd_rental/dvds/model"
	. "github.com/sub0zero/go-sqlbuilder/tests/.test_files/dvd_rental/dvds/table"
	model2 "github.com/sub0zero/go-sqlbuilder/tests/.test_files/dvd_rental/test_sample/model"
	. "github.com/sub0zero/go-sqlbuilder/tests/.test_files/dvd_rental/test_sample/table"
	"gotest.tools/assert"
	"testing"
)

func TestSelect_ScanToStruct(t *testing.T) {
	expectedSql := `
SELECT actor.actor_id AS "actor.actor_id",
     actor.first_name AS "actor.first_name",
     actor.last_name AS "actor.last_name",
     actor.last_update AS "actor.last_update"
FROM dvds.actor
WHERE actor.actor_id = 1;
`

	query := Actor.
		SELECT(Actor.AllColumns).
		WHERE(Actor.ActorID.EqL(1))

	assertQuery(t, query, expectedSql, 1)

	actor := model.Actor{}
	err := query.Query(db, &actor)

	assert.NilError(t, err)

	expectedActor := model.Actor{
		ActorID:    1,
		FirstName:  "Penelope",
		LastName:   "Guiness",
		LastUpdate: *timeWithoutTimeZone("2013-05-26 14:47:57.62", 2),
	}

	assert.DeepEqual(t, actor, expectedActor)
}

func TestClassicSelect(t *testing.T) {
	expectedSql := `
SELECT payment.payment_id AS "payment.payment_id",
     payment.customer_id AS "payment.customer_id",
     payment.staff_id AS "payment.staff_id",
     payment.rental_id AS "payment.rental_id",
     payment.amount AS "payment.amount",
     payment.payment_date AS "payment.payment_date",
     customer.customer_id AS "customer.customer_id",
     customer.store_id AS "customer.store_id",
     customer.first_name AS "customer.first_name",
     customer.last_name AS "customer.last_name",
     customer.email AS "customer.email",
     customer.address_id AS "customer.address_id",
     customer.activebool AS "customer.activebool",
     customer.create_date AS "customer.create_date",
     customer.last_update AS "customer.last_update",
     customer.active AS "customer.active"
FROM dvds.payment
     JOIN dvds.customer ON payment.customer_id = customer.customer_id
ORDER BY payment.payment_id ASC
LIMIT 30;
`

	query := SELECT(Payment.AllColumns, Customer.AllColumns).
		FROM(Payment.INNER_JOIN(Customer, Payment.CustomerID.Eq(Customer.CustomerID))).
		ORDER_BY(Payment.PaymentID.ASC()).
		LIMIT(30)

	assertQuery(t, query, expectedSql, int64(30))

	dest := []model.Payment{}

	err := query.Query(db, &dest)

	assert.NilError(t, err)
	assert.Equal(t, len(dest), 30)
}

func TestSelect_ScanToSlice(t *testing.T) {
	expectedSql := `
SELECT customer.customer_id AS "customer.customer_id",
     customer.store_id AS "customer.store_id",
     customer.first_name AS "customer.first_name",
     customer.last_name AS "customer.last_name",
     customer.email AS "customer.email",
     customer.address_id AS "customer.address_id",
     customer.activebool AS "customer.activebool",
     customer.create_date AS "customer.create_date",
     customer.last_update AS "customer.last_update",
     customer.active AS "customer.active"
FROM dvds.customer
ORDER BY customer.customer_id ASC;
`
	customers := []model.Customer{}

	query := Customer.SELECT(Customer.AllColumns).ORDER_BY(Customer.CustomerID.ASC())

	assertQuery(t, query, expectedSql)

	err := query.Query(db, &customers)
	assert.NilError(t, err)

	assert.Equal(t, len(customers), 599)

	assert.DeepEqual(t, customer0, customers[0])
	assert.DeepEqual(t, customer1, customers[1])
	assert.DeepEqual(t, lastCustomer, customers[598])
}

func TestSelectAndUnionInProjection(t *testing.T) {
	expectedSql := `
SELECT payment.payment_id AS "payment.payment_id",
     (
          SELECT customer.customer_id AS "customer.customer_id"
          FROM dvds.customer
          LIMIT 1
     ),
     (
          (
               (
                    SELECT payment.payment_id AS "payment.payment_id"
                    FROM dvds.payment
                    LIMIT 1
                    OFFSET 10
               )
               UNION
               (
                    SELECT payment.payment_id AS "payment.payment_id"
                    FROM dvds.payment
                    LIMIT 1
                    OFFSET 2
               )
          )
          LIMIT 1
     )
FROM dvds.payment
LIMIT 12;
`

	query := Payment.
		SELECT(
			Payment.PaymentID,
			Customer.SELECT(Customer.CustomerID).LIMIT(1),
			UNION(
				Payment.SELECT(Payment.PaymentID).LIMIT(1).OFFSET(10),
				Payment.SELECT(Payment.PaymentID).LIMIT(1).OFFSET(2),
			).LIMIT(1),
		).
		LIMIT(12)

	fmt.Println(query.Sql())
	assertQuery(t, query, expectedSql, int64(1), int64(1), int64(10), int64(1), int64(2), int64(1), int64(12))
}

func TestJoinQueryStruct(t *testing.T) {

	expectedSql := `
SELECT film_actor.actor_id AS "film_actor.actor_id",
     film_actor.film_id AS "film_actor.film_id",
     film_actor.last_update AS "film_actor.last_update",
     film.film_id AS "film.film_id",
     film.title AS "film.title",
     film.description AS "film.description",
     film.release_year AS "film.release_year",
     film.language_id AS "film.language_id",
     film.rental_duration AS "film.rental_duration",
     film.rental_rate AS "film.rental_rate",
     film.length AS "film.length",
     film.replacement_cost AS "film.replacement_cost",
     film.rating AS "film.rating",
     film.last_update AS "film.last_update",
     film.special_features AS "film.special_features",
     film.fulltext AS "film.fulltext",
     language.language_id AS "language.language_id",
     language.name AS "language.name",
     language.last_update AS "language.last_update",
     actor.actor_id AS "actor.actor_id",
     actor.first_name AS "actor.first_name",
     actor.last_name AS "actor.last_name",
     actor.last_update AS "actor.last_update",
     inventory.inventory_id AS "inventory.inventory_id",
     inventory.film_id AS "inventory.film_id",
     inventory.store_id AS "inventory.store_id",
     inventory.last_update AS "inventory.last_update",
     rental.rental_id AS "rental.rental_id",
     rental.rental_date AS "rental.rental_date",
     rental.inventory_id AS "rental.inventory_id",
     rental.customer_id AS "rental.customer_id",
     rental.return_date AS "rental.return_date",
     rental.staff_id AS "rental.staff_id",
     rental.last_update AS "rental.last_update"
FROM dvds.film_actor
     JOIN dvds.actor ON film_actor.actor_id = actor.actor_id
     JOIN dvds.film ON film_actor.film_id = film.film_id
     JOIN dvds.language ON film.language_id = language.language_id
     JOIN dvds.inventory ON inventory.film_id = film.film_id
     JOIN dvds.rental ON rental.inventory_id = inventory.inventory_id
ORDER BY film.film_id ASC
LIMIT 50;
`
	for i := 0; i < 1; i++ {
		query := FilmActor.
			INNER_JOIN(Actor, FilmActor.ActorID.Eq(Actor.ActorID)).
			INNER_JOIN(Film, FilmActor.FilmID.Eq(Film.FilmID)).
			INNER_JOIN(Language, Film.LanguageID.Eq(Language.LanguageID)).
			INNER_JOIN(Inventory, Inventory.FilmID.Eq(Film.FilmID)).
			INNER_JOIN(Rental, Rental.InventoryID.Eq(Inventory.InventoryID)).
			SELECT(
				FilmActor.AllColumns,
				Film.AllColumns,
				Language.AllColumns,
				Actor.AllColumns,
				Inventory.AllColumns,
				Rental.AllColumns,
			).
			//WHERE(FilmActor.ActorID.GtEqL(1).AND(FilmActor.ActorID.LtEqL(2))).
			ORDER_BY(Film.FilmID.ASC()).
			LIMIT(50)

		assertQuery(t, query, expectedSql, int64(50))

		var languageActorFilm []struct {
			model.Language

			Films []struct {
				model.Film
				Actors []struct {
					model.Actor
				}

				Inventory []struct {
					model.Inventory

					Rental []model.Rental
				}
			}
		}

		err := query.Query(db, &languageActorFilm)

		assert.NilError(t, err)
		assert.Equal(t, len(languageActorFilm), 1)
		assert.Equal(t, len(languageActorFilm[0].Films), 1)
		assert.Equal(t, len(languageActorFilm[0].Films[0].Actors), 10)
	}

}

func TestJoinQuerySlice(t *testing.T) {
	expectedSql := `
SELECT language.language_id AS "language.language_id",
     language.name AS "language.name",
     language.last_update AS "language.last_update",
     film.film_id AS "film.film_id",
     film.title AS "film.title",
     film.description AS "film.description",
     film.release_year AS "film.release_year",
     film.language_id AS "film.language_id",
     film.rental_duration AS "film.rental_duration",
     film.rental_rate AS "film.rental_rate",
     film.length AS "film.length",
     film.replacement_cost AS "film.replacement_cost",
     film.rating AS "film.rating",
     film.last_update AS "film.last_update",
     film.special_features AS "film.special_features",
     film.fulltext AS "film.fulltext"
FROM dvds.film
     JOIN dvds.language ON film.language_id = language.language_id
WHERE film.rating = 'NC-17'
LIMIT 15;
`

	type FilmsPerLanguage struct {
		Language *model.Language
		Film     []model.Film
	}

	filmsPerLanguage := []FilmsPerLanguage{}
	limit := 15

	query := Film.
		INNER_JOIN(Language, Film.LanguageID.Eq(Language.LanguageID)).
		SELECT(Language.AllColumns, Film.AllColumns).
		WHERE(Film.Rating.EqString(model.MpaaRating_NC17.String())).
		LIMIT(15)

	assertQuery(t, query, expectedSql, model.MpaaRating_NC17.String(), int64(15))

	err := query.Query(db, &filmsPerLanguage)

	assert.NilError(t, err)
	assert.Equal(t, len(filmsPerLanguage), 1)
	assert.Equal(t, len(filmsPerLanguage[0].Film), limit)

	englishFilms := filmsPerLanguage[0]

	assert.Equal(t, *englishFilms.Film[0].Rating, model.MpaaRating_NC17)

	filmsPerLanguageWithPtrs := []*FilmsPerLanguage{}
	err = query.Query(db, &filmsPerLanguageWithPtrs)

	assert.NilError(t, err)
	assert.Equal(t, len(filmsPerLanguage), 1)
	assert.Equal(t, len(filmsPerLanguage[0].Film), limit)
}

func TestJoinQuerySliceWithPtrs(t *testing.T) {
	type FilmsPerLanguage struct {
		Language model.Language
		Film     *[]*model.Film
	}

	limit := int64(3)

	query := Film.INNER_JOIN(Language, Film.LanguageID.Eq(Language.LanguageID)).
		SELECT(Language.AllColumns, Film.AllColumns).
		LIMIT(limit)

	filmsPerLanguageWithPtrs := []*FilmsPerLanguage{}
	err := query.Query(db, &filmsPerLanguageWithPtrs)

	assert.NilError(t, err)
	assert.Equal(t, len(filmsPerLanguageWithPtrs), 1)
	assert.Equal(t, len(*filmsPerLanguageWithPtrs[0].Film), int(limit))
}

func TestSelect_WithoutUniqueColumnSelected(t *testing.T) {
	query := Customer.SELECT(Customer.FirstName, Customer.LastName, Customer.Email)

	customers := []model.Customer{}

	err := query.Query(db, &customers)

	assert.NilError(t, err)

	//spew.Dump(customers)

	assert.Equal(t, len(customers), 599)
}

func TestSelectOrderByAscDesc(t *testing.T) {
	customersAsc := []model.Customer{}

	err := Customer.SELECT(Customer.CustomerID, Customer.FirstName, Customer.LastName).
		ORDER_BY(Customer.FirstName.ASC()).
		Query(db, &customersAsc)

	assert.NilError(t, err)

	firstCustomerAsc := customersAsc[0]
	lastCustomerAsc := customersAsc[len(customersAsc)-1]

	customersDesc := []model.Customer{}
	err = Customer.SELECT(Customer.CustomerID, Customer.FirstName, Customer.LastName).
		ORDER_BY(Customer.FirstName.DESC()).
		Query(db, &customersDesc)

	assert.NilError(t, err)

	firstCustomerDesc := customersDesc[0]
	lastCustomerDesc := customersDesc[len(customersAsc)-1]

	assert.DeepEqual(t, firstCustomerAsc, lastCustomerDesc)
	assert.DeepEqual(t, lastCustomerAsc, firstCustomerDesc)

	customersAscDesc := []model.Customer{}
	err = Customer.SELECT(Customer.CustomerID, Customer.FirstName, Customer.LastName).
		ORDER_BY(Customer.FirstName.ASC(), Customer.LastName.DESC()).
		Query(db, &customersAscDesc)

	assert.NilError(t, err)

	customerAscDesc326 := model.Customer{
		CustomerID: 67,
		FirstName:  "Kelly",
		LastName:   "Torres",
	}

	customerAscDesc327 := model.Customer{
		CustomerID: 546,
		FirstName:  "Kelly",
		LastName:   "Knott",
	}

	assert.DeepEqual(t, customerAscDesc326, customersAscDesc[326])
	assert.DeepEqual(t, customerAscDesc327, customersAscDesc[327])
}

func TestSelectFullJoin(t *testing.T) {
	expectedSql := `
SELECT customer.customer_id AS "customer.customer_id",
     customer.store_id AS "customer.store_id",
     customer.first_name AS "customer.first_name",
     customer.last_name AS "customer.last_name",
     customer.email AS "customer.email",
     customer.address_id AS "customer.address_id",
     customer.activebool AS "customer.activebool",
     customer.create_date AS "customer.create_date",
     customer.last_update AS "customer.last_update",
     customer.active AS "customer.active",
     address.address_id AS "address.address_id",
     address.address AS "address.address",
     address.address2 AS "address.address2",
     address.district AS "address.district",
     address.city_id AS "address.city_id",
     address.postal_code AS "address.postal_code",
     address.phone AS "address.phone",
     address.last_update AS "address.last_update"
FROM dvds.customer
     FULL JOIN dvds.address ON customer.address_id = address.address_id
ORDER BY customer.customer_id ASC;
`
	query := Customer.
		FULL_JOIN(Address, Customer.AddressID.Eq(Address.AddressID)).
		SELECT(Customer.AllColumns, Address.AllColumns).
		ORDER_BY(Customer.CustomerID.ASC())

	assertQuery(t, query, expectedSql)

	allCustomersAndAddress := []struct {
		Address  *model.Address
		Customer *model.Customer
	}{}

	err := query.Query(db, &allCustomersAndAddress)

	assert.NilError(t, err)
	assert.Equal(t, len(allCustomersAndAddress), 603)

	assert.DeepEqual(t, allCustomersAndAddress[0].Customer, &customer0)
	assert.Assert(t, allCustomersAndAddress[0].Address != nil)

	lastCustomerAddress := allCustomersAndAddress[len(allCustomersAndAddress)-1]

	assert.Assert(t, lastCustomerAddress.Customer == nil)
	assert.Assert(t, lastCustomerAddress.Address != nil)

}

func TestSelectFullCrossJoin(t *testing.T) {
	expectedSql := `
SELECT customer.customer_id AS "customer.customer_id",
     customer.store_id AS "customer.store_id",
     customer.first_name AS "customer.first_name",
     customer.last_name AS "customer.last_name",
     customer.email AS "customer.email",
     customer.address_id AS "customer.address_id",
     customer.activebool AS "customer.activebool",
     customer.create_date AS "customer.create_date",
     customer.last_update AS "customer.last_update",
     customer.active AS "customer.active",
     address.address_id AS "address.address_id",
     address.address AS "address.address",
     address.address2 AS "address.address2",
     address.district AS "address.district",
     address.city_id AS "address.city_id",
     address.postal_code AS "address.postal_code",
     address.phone AS "address.phone",
     address.last_update AS "address.last_update"
FROM dvds.customer
     CROSS JOIN dvds.address
ORDER BY customer.customer_id ASC
LIMIT 1000;
`
	query := Customer.
		CROSS_JOIN(Address).
		SELECT(Customer.AllColumns, Address.AllColumns).
		ORDER_BY(Customer.CustomerID.ASC()).
		LIMIT(1000)

	assertQuery(t, query, expectedSql, int64(1000))

	var customerAddresCrosJoined []struct {
		model.Customer
		model.Address
	}

	err := query.Query(db, &customerAddresCrosJoined)

	assert.Equal(t, len(customerAddresCrosJoined), 1000)

	assert.NilError(t, err)
}

func TestSelecSelfJoin1(t *testing.T) {

	var expectedSql = `
SELECT employee.employee_id AS "employee.employee_id",
     employee.first_name AS "employee.first_name",
     employee.last_name AS "employee.last_name",
     employee.manager_id AS "employee.manager_id",
     manager.employee_id AS "manager.employee_id",
     manager.first_name AS "manager.first_name",
     manager.last_name AS "manager.last_name",
     manager.manager_id AS "manager.manager_id"
FROM test_sample.employee
     LEFT JOIN test_sample.employee AS manager ON manager.employee_id = employee.manager_id
ORDER BY employee.employee_id;
`

	manager := Employee.AS("manager")
	query := Employee.
		LEFT_JOIN(manager, manager.EmployeeID.Eq(Employee.ManagerID)).
		SELECT(Employee.AllColumns, manager.AllColumns).
		ORDER_BY(Employee.EmployeeID)

	assertQuery(t, query, expectedSql)

	var dest []struct {
		model2.Employee

		Manager *model2.Employee
	}

	err := query.Query(db, &dest)

	assert.NilError(t, err)
	assert.Equal(t, len(dest), 8)
	assert.DeepEqual(t, dest[0].Employee, model2.Employee{
		EmployeeID: 1,
		FirstName:  "Windy",
		LastName:   "Hays",
		ManagerID:  nil,
	})

	assert.Assert(t, dest[0].Manager == nil)

	assert.DeepEqual(t, dest[7].Employee, model2.Employee{
		EmployeeID: 8,
		FirstName:  "Salley",
		LastName:   "Lester",
		ManagerID:  int32Ptr(3),
	})
}

func TestSelectSelfJoin(t *testing.T) {
	expectedSql := `
SELECT f1.film_id AS "f1.film_id",
     f1.title AS "f1.title",
     f1.description AS "f1.description",
     f1.release_year AS "f1.release_year",
     f1.language_id AS "f1.language_id",
     f1.rental_duration AS "f1.rental_duration",
     f1.rental_rate AS "f1.rental_rate",
     f1.length AS "f1.length",
     f1.replacement_cost AS "f1.replacement_cost",
     f1.rating AS "f1.rating",
     f1.last_update AS "f1.last_update",
     f1.special_features AS "f1.special_features",
     f1.fulltext AS "f1.fulltext",
     f2.film_id AS "f2.film_id",
     f2.title AS "f2.title",
     f2.description AS "f2.description",
     f2.release_year AS "f2.release_year",
     f2.language_id AS "f2.language_id",
     f2.rental_duration AS "f2.rental_duration",
     f2.rental_rate AS "f2.rental_rate",
     f2.length AS "f2.length",
     f2.replacement_cost AS "f2.replacement_cost",
     f2.rating AS "f2.rating",
     f2.last_update AS "f2.last_update",
     f2.special_features AS "f2.special_features",
     f2.fulltext AS "f2.fulltext"
FROM dvds.film AS f1
     JOIN dvds.film AS f2 ON (f1.film_id < f2.film_id AND f1.length = f2.length)
ORDER BY f1.film_id ASC;
`
	f1 := Film.AS("f1")

	f2 := Film.AS("f2")

	query := f1.
		INNER_JOIN(f2, f1.FilmID.Lt(f2.FilmID).AND(f1.Length.Eq(f2.Length))).
		SELECT(f1.AllColumns, f2.AllColumns).
		ORDER_BY(f1.FilmID.ASC())

	assertQuery(t, query, expectedSql)

	type F1 model.Film
	type F2 model.Film

	theSameLengthFilms := []struct {
		F1 F1
		F2 F2
	}{}

	err := query.Query(db, &theSameLengthFilms)

	assert.NilError(t, err)

	//spew.Dump(theSameLengthFilms)

	//assert.Equal(t, len(theSameLengthFilms), 100)
}

func TestSelectAliasColumn(t *testing.T) {
	expectedSql := `
SELECT f1.title AS "thesame_length_films.title1",
     f2.title AS "thesame_length_films.title2",
     f1.length AS "thesame_length_films.length"
FROM dvds.film AS f1
     JOIN dvds.film AS f2 ON (f1.film_id != f2.film_id AND f1.length = f2.length)
ORDER BY f1.length ASC, f1.title ASC, f2.title ASC
LIMIT 1000;
`
	f1 := Film.AS("f1")
	f2 := Film.AS("f2")

	query := f1.
		INNER_JOIN(f2, f1.FilmID.NotEq(f2.FilmID).AND(f1.Length.Eq(f2.Length))).
		SELECT(f1.Title.AS("thesame_length_films.title1"),
			f2.Title.AS("thesame_length_films.title2"),
			f1.Length.AS("thesame_length_films.length")).
		ORDER_BY(f1.Length.ASC(), f1.Title.ASC(), f2.Title.ASC()).
		LIMIT(1000)

	assertQuery(t, query, expectedSql, int64(1000))

	type thesameLengthFilms struct {
		Title1 string
		Title2 string
		Length int16
	}
	films := []thesameLengthFilms{}

	err := query.Query(db, &films)

	assert.NilError(t, err)

	//spew.Dump(films)

	assert.Equal(t, len(films), 1000)
	assert.DeepEqual(t, films[0], thesameLengthFilms{"Alien Center", "Iron Moon", 46})
}

//
//type Manager staff
//
//type staff struct {
//	StaffID   int32 `sql:"unique"`
//	FirstName string
//	LastName  string
//	//Address    *model.Address
//	//Email      *string
//	//StoreID    int16
//	//Active     bool
//	//Username   string
//	//Password   *string
//	//LastUpdate time.Time
//	*Manager //`sqlbuilder:"manager"`
//}
//
//func TestSelectSelfReferenceType(t *testing.T) {
//
//	expectedSql := `
//SELECT DISTINCT staff.staff_id AS "staff.staff_id",
//     staff.first_name AS "staff.first_name",
//     staff.last_name AS "staff.last_name",
//     address.address_id AS "address.address_id",
//     address.address AS "address.address",
//     address.address2 AS "address.address2",
//     address.district AS "address.district",
//     address.city_id AS "address.city_id",
//     address.postal_code AS "address.postal_code",
//     address.phone AS "address.phone",
//     address.last_update AS "address.last_update",
//     manager.staff_id AS "manager.staff_id",
//     manager.first_name AS "manager.first_name"
//FROM dvds.staff
//     JOIN dvds.address ON staff.address_id = address.address_id
//     JOIN dvds.staff AS manager ON staff.staff_id = manager.staff_id;
//`
//	manager := Staff.AS("manager")
//
//	query := Staff.
//		INNER_JOIN(Address, Staff.AddressID.Eq(Address.AddressID)).
//		INNER_JOIN(manager, Staff.StaffID.Eq(manager.StaffID)).
//		SELECT(Staff.StaffID, Staff.FirstName, Staff.LastName, Address.AllColumns, manager.StaffID, manager.FirstName).
//		DISTINCT()
//
//	assertQuery(t, query, expectedSql)
//
//	staffs := []staff{}
//
//	err := query.Query(db, &staffs)
//
//	assert.NilError(t, err)
//
//	fmt.Println(query.DebugSql())
//	//spew.Dump(staffs)
//}

func TestSubQuery(t *testing.T) {

	//selectStmtTable := Actor.SELECT(Actor.FirstName, Actor.LastName).AsTable("table_expression")
	//
	//query := selectStmtTable.SELECT(
	//	selectStmtTable.RefStringColumn(Actor.FirstName).AS("nesto"),
	//	selectStmtTable.RefIntColumnName("actor.last_name").AS("nesto2"),
	//	)
	//
	//queryStr, args, err := query.Sql()
	//
	//assert.NilError(t, err)
	//
	//fmt.Println(queryStr)
	//
	//avrgCustomer := NumExp(Customer.SELECT(Customer.LastName).LIMIT(1))
	//
	//Customer.
	//	INNER_JOIN(selectStmtTable, Customer.LastName.Eq(selectStmtTable.RefStringColumn(Actor.FirstName))).
	//	SELECT(Customer.AllColumns, selectStmtTable.RefIntColumnName("first_name")).
	//	WHERE(Actor.LastName.Neq(avrgCustomer))

	expectedQuery := `
SELECT actor.actor_id AS "actor.actor_id",
     actor.first_name AS "actor.first_name",
     actor.last_name AS "actor.last_name",
     actor.last_update AS "actor.last_update",
     film_actor.actor_id AS "film_actor.actor_id",
     film_actor.film_id AS "film_actor.film_id",
     film_actor.last_update AS "film_actor.last_update",
     films."film.title" AS "film.title",
     films."film.rating" AS "film.rating"
FROM dvds.actor
     JOIN dvds.film_actor ON actor.actor_id = film_actor.film_id
     JOIN (
          SELECT film.film_id AS "film.film_id",
               film.title AS "film.title",
               film.rating AS "film.rating"
          FROM dvds.film
          WHERE film.rating = 'R'
     ) AS films ON film_actor.film_id = films."film.film_id";
`

	rFilmsOnly := Film.SELECT(Film.FilmID, Film.Title, Film.Rating).
		WHERE(Film.Rating.EqString("R")).
		AsTable("films")

	query := Actor.INNER_JOIN(FilmActor, Actor.ActorID.Eq(FilmActor.FilmID)).
		INNER_JOIN(rFilmsOnly, FilmActor.FilmID.Eq(rFilmsOnly.RefIntColumn(Film.FilmID))).
		SELECT(
			Actor.AllColumns,
			FilmActor.AllColumns,
			rFilmsOnly.RefStringColumn(Film.Title).AS("film.title"),
			rFilmsOnly.RefStringColumn(Film.Rating).AS("film.rating"),
		)

	fmt.Println(query.Sql())

	assertQuery(t, query, expectedQuery, "R")

	dest := []model.Actor{}

	err := query.Query(db, &dest)

	assert.NilError(t, err)
}

func TestSelectFunctions(t *testing.T) {
	expectedQuery := `
SELECT MAX(film.rental_rate) AS "max_film_rate"
FROM dvds.film;
`
	query := Film.SELECT(MAX(Film.RentalRate).AS("max_film_rate"))

	assertQuery(t, query, expectedQuery)

	ret := struct {
		MaxFilmRate float64
	}{}

	err := query.Query(db, &ret)

	assert.NilError(t, err)
	assert.Equal(t, ret.MaxFilmRate, 4.99)
}

func TestSelectQueryScalar(t *testing.T) {
	expectedSql := `
SELECT film.film_id AS "film.film_id",
     film.title AS "film.title",
     film.description AS "film.description",
     film.release_year AS "film.release_year",
     film.language_id AS "film.language_id",
     film.rental_duration AS "film.rental_duration",
     film.rental_rate AS "film.rental_rate",
     film.length AS "film.length",
     film.replacement_cost AS "film.replacement_cost",
     film.rating AS "film.rating",
     film.last_update AS "film.last_update",
     film.special_features AS "film.special_features",
     film.fulltext AS "film.fulltext"
FROM dvds.film
WHERE film.rental_rate = (
          SELECT MAX(film.rental_rate)
          FROM dvds.film
     )
ORDER BY film.film_id ASC;
`

	maxFilmRentalRate := NumExp(Film.SELECT(MAX(Film.RentalRate)))

	query := Film.
		SELECT(Film.AllColumns).
		WHERE(Film.RentalRate.Eq(maxFilmRentalRate)).
		ORDER_BY(Film.FilmID.ASC())

	fmt.Println(query.Sql())
	assertQuery(t, query, expectedSql)

	maxRentalRateFilms := []model.Film{}
	err := query.Query(db, &maxRentalRateFilms)

	assert.NilError(t, err)

	assert.Equal(t, len(maxRentalRateFilms), 336)

	gRating := model.MpaaRating_G

	assert.DeepEqual(t, maxRentalRateFilms[0], model.Film{
		FilmID:          2,
		Title:           "Ace Goldfinger",
		Description:     stringPtr("A Astounding Epistle of a Database Administrator And a Explorer who must Find a Car in Ancient China"),
		ReleaseYear:     int32Ptr(2006),
		LanguageID:      1,
		RentalRate:      4.99,
		Length:          int16Ptr(48),
		ReplacementCost: 12.99,
		Rating:          &gRating,
		RentalDuration:  3,
		LastUpdate:      *timeWithoutTimeZone("2013-05-26 14:50:58.951", 3),
		SpecialFeatures: stringPtr("{Trailers,\"Deleted Scenes\"}"),
		Fulltext:        "'ace':1 'administr':9 'ancient':19 'astound':4 'car':17 'china':20 'databas':8 'epistl':5 'explor':12 'find':15 'goldfing':2 'must':14",
	})
}

func TestSelectGroupByHaving(t *testing.T) {
	expectedSql := `
SELECT payment.customer_id AS "customer_payment_sum.customer_id",
     SUM(payment.amount) AS "customer_payment_sum.amount_sum"
FROM dvds.payment
GROUP BY payment.customer_id
HAVING SUM(payment.amount) > 100
ORDER BY SUM(payment.amount) ASC;
`
	customersPaymentQuery := Payment.
		SELECT(
			Payment.CustomerID.AS("customer_payment_sum.customer_id"),
			SUM(Payment.Amount).AS("customer_payment_sum.amount_sum"),
		).
		GROUP_BY(Payment.CustomerID).
		ORDER_BY(SUM(Payment.Amount).ASC()).
		HAVING(SUM(Payment.Amount).Gt(NewNumericLiteral(100)))

	assertQuery(t, customersPaymentQuery, expectedSql, 100)

	type CustomerPaymentSum struct {
		CustomerID int16
		AmountSum  float64
	}

	customerPaymentSum := []CustomerPaymentSum{}

	err := customersPaymentQuery.Query(db, &customerPaymentSum)

	assert.NilError(t, err)

	assert.Equal(t, len(customerPaymentSum), 296)
	assert.DeepEqual(t, customerPaymentSum[0], CustomerPaymentSum{
		CustomerID: 135,
		AmountSum:  100.72,
	})
}

func TestSelectGroupBy2(t *testing.T) {
	expectedSql := `
SELECT customer.customer_id AS "customer.customer_id",
     customer.store_id AS "customer.store_id",
     customer.first_name AS "customer.first_name",
     customer.last_name AS "customer.last_name",
     customer.email AS "customer.email",
     customer.address_id AS "customer.address_id",
     customer.activebool AS "customer.activebool",
     customer.create_date AS "customer.create_date",
     customer.last_update AS "customer.last_update",
     customer.active AS "customer.active",
     customer_payment_sum.amount_sum AS "customer_with_amounts.amount_sum"
FROM dvds.customer
     JOIN (
          SELECT payment.customer_id AS "payment.customer_id",
               SUM(payment.amount) AS "amount_sum"
          FROM dvds.payment
          GROUP BY payment.customer_id
     ) AS customer_payment_sum ON customer.customer_id = customer_payment_sum."payment.customer_id"
ORDER BY customer_payment_sum.amount_sum ASC;
`

	customersPaymentSubQuery := Payment.
		SELECT(
			Payment.CustomerID,
			SUM(Payment.Amount).AS("amount_sum"),
		).
		GROUP_BY(Payment.CustomerID)

	customersPaymentTable := customersPaymentSubQuery.AsTable("customer_payment_sum")
	amountSumColumn := customersPaymentTable.RefIntColumnName("amount_sum")

	query := Customer.
		INNER_JOIN(customersPaymentTable, Customer.CustomerID.Eq(customersPaymentTable.RefIntColumn(Payment.CustomerID))).
		SELECT(Customer.AllColumns, amountSumColumn.AS("customer_with_amounts.amount_sum")).
		ORDER_BY(amountSumColumn.ASC())

	assertQuery(t, query, expectedSql)

	type CustomerWithAmounts struct {
		Customer  *model.Customer
		AmountSum float64
	}
	customersWithAmounts := []CustomerWithAmounts{}

	err := query.Query(db, &customersWithAmounts)
	assert.NilError(t, err)
	assert.Equal(t, len(customersWithAmounts), 599)

	assert.DeepEqual(t, customersWithAmounts[0].Customer, &model.Customer{
		CustomerID: 318,
		StoreID:    1,
		FirstName:  "Brian",
		LastName:   "Wyman",
		AddressID:  323,
		Email:      stringPtr("brian.wyman@sakilacustomer.org"),
		Activebool: true,
		CreateDate: *timeWithoutTimeZone("2006-02-14 00:00:00", 0),
		LastUpdate: timeWithoutTimeZone("2013-05-26 14:49:45.738", 3),
		Active:     int32Ptr(1),
	})

	assert.Equal(t, customersWithAmounts[0].AmountSum, 27.93)
}

func TestSelectTimeColumns(t *testing.T) {

	expectedSql := `
SELECT payment.payment_id AS "payment.payment_id",
     payment.customer_id AS "payment.customer_id",
     payment.staff_id AS "payment.staff_id",
     payment.rental_id AS "payment.rental_id",
     payment.amount AS "payment.amount",
     payment.payment_date AS "payment.payment_date"
FROM dvds.payment
WHERE payment.payment_date <= '2007-02-14 22:16:01'
ORDER BY payment.payment_date ASC;
`

	query := Payment.SELECT(Payment.AllColumns).
		WHERE(Payment.PaymentDate.LtEqL("2007-02-14 22:16:01")).
		ORDER_BY(Payment.PaymentDate.ASC())

	assertQuery(t, query, expectedSql, "2007-02-14 22:16:01")

	payments := []model.Payment{}

	err := query.Query(db, &payments)

	assert.NilError(t, err)

	//spew.Dump(payments)

	assert.Equal(t, len(payments), 9)
	assert.DeepEqual(t, payments[0], model.Payment{
		PaymentID:   17793,
		CustomerID:  416,
		StaffID:     2,
		RentalID:    1158,
		Amount:      2.99,
		PaymentDate: *timeWithoutTimeZone("2007-02-14 21:21:59.996577", 6),
	})
}

func TestUnion(t *testing.T) {
	expectedQuery := `
(
     (
          SELECT payment.payment_id AS "payment.payment_id",
               payment.amount AS "payment.amount"
          FROM dvds.payment
          WHERE payment.amount <= 100
     )
     UNION ALL
     (
          SELECT payment.payment_id AS "payment.payment_id",
               payment.amount AS "payment.amount"
          FROM dvds.payment
          WHERE payment.amount >= 200
     )
)
ORDER BY "payment.payment_id" ASC, "payment.amount" DESC
LIMIT 10
OFFSET 20;
`
	query := UNION_ALL(
		Payment.
			SELECT(Payment.PaymentID.AS("payment.payment_id"), Payment.Amount).
			WHERE(Payment.Amount.LtEqL(100)),
		Payment.
			SELECT(Payment.PaymentID, Payment.Amount).
			WHERE(Payment.Amount.GtEqL(200)),
	).
		ORDER_BY(RefColumn("payment.payment_id").ASC(), Payment.Amount.DESC()).
		LIMIT(10).
		OFFSET(20)

	queryStr, _, _ := query.Sql()

	fmt.Println("-" + queryStr + "-")
	assertQuery(t, query, expectedQuery, int(100), int(200), int64(10), int64(20))

	dest := []model.Payment{}

	err := query.Query(db, &dest)

	assert.NilError(t, err)
	assert.Equal(t, len(dest), 10)
	assert.DeepEqual(t, dest[0], model.Payment{
		PaymentID: 17523,
		Amount:    4.99,
	})
	assert.DeepEqual(t, dest[1], model.Payment{
		PaymentID: 17524,
		Amount:    0.99,
	})
	assert.DeepEqual(t, dest[9], model.Payment{
		PaymentID: 17532,
		Amount:    8.99,
	})
}

func TestSelectWithCase(t *testing.T) {
	expectedQuery := `
SELECT (CASE payment.staff_id WHEN 1 THEN 'ONE' WHEN 2 THEN 'TWO' WHEN 3 THEN 'THREE' ELSE 'OTHER' END) AS "staff_id_num"
FROM dvds.payment
ORDER BY payment.payment_id ASC
LIMIT 20;
`
	query := Payment.SELECT(
		CASE(Payment.StaffID).
			WHEN(IntLiteral(1)).THEN(Literal("ONE")).
			WHEN(IntLiteral(2)).THEN(Literal("TWO")).
			WHEN(IntLiteral(3)).THEN(Literal("THREE")).
			ELSE(Literal("OTHER")).AS("staff_id_num"),
	).
		ORDER_BY(Payment.PaymentID.ASC()).
		LIMIT(20)

	assertQuery(t, query, expectedQuery, 1, "ONE", 2, "TWO", 3, "THREE", "OTHER", int64(20))

	dest := []struct {
		StaffIdNum string
	}{}

	err := query.Query(db, &dest)

	assert.NilError(t, err)
	assert.Equal(t, len(dest), 20)
	assert.Equal(t, dest[0].StaffIdNum, "TWO")
	assert.Equal(t, dest[1].StaffIdNum, "ONE")
}

func TestLockTable(t *testing.T) {
	expectedSql := `
LOCK TABLE dvds.address IN EXCLUSIVE MODE NOWAIT;
`
	query := Address.LOCK().IN(LOCK_EXCLUSIVE).NOWAIT()

	querySql, _, _ := query.Sql()
	fmt.Println("-" + querySql + "-")

	assertQuery(t, query, expectedSql)

	tx, _ := db.Begin()

	_, err := query.Execute(tx)

	assert.NilError(t, err)
}
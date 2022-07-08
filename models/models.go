package models 

import (
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"os"
	"fmt"
	_ "strings"
	"log"
	"time"
)

var (
	DebugLogger *log.Logger 
	InfoLogger *log.Logger 
	ErrorLogger *log.Logger 
	WarnLogger *log.Logger 
)
var (
	POSTGRES_USER = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DATABASE = os.Getenv("POSTGRES_DATABASE")
	POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
)

var (
	Database, error = gorm.Open(postgres.New(
		postgres.Config{
			DSN: fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		POSTGRES_HOST, POSTGRES_PORT, POSTGRES_DATABASE, POSTGRES_USER, POSTGRES_PASSWORD),
			PreferSimpleProtocol: true,
		},
	))
)

func init() {
	LogFile, error := os.OpenFile("databaseLogs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if error != nil {panic("Failed to Initialize Database Log File.")}

	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(LogFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
}

type Product struct {
	gorm.Model 
	ProductName string `gorm:"VARCHAR(100) NOT NULL"`
	ProductDescription string `gorm:"VARCHAR(100) NOT NULL DEFAULT 'This Product Has No Description'"`
	ProductPrice string `gorm:"NUMERIC(10, 5) NOT NULL"`
	Currency string `gorm:"VARCHAR(10) NOT NULL"`
}

type Customer struct {
	gorm.Model 
	Username string `gorm:"VARCHAR(100) NOT NULL"`
	Password string `gorm:"VARCHAR(100) NOT NULL"`
	Email string `gorm:"VARCHAR(100) NOT NULL"`
	PurchasedProducts Product  `gorm:"foreignKey:Product;references:ProductId"`
	CreatedAt time.Time `gorm:"DATE DEFAULT CURRENT DATE"`
}

type Cart struct {
	gorm.Model 
	Owner Customer `gorm:"foreignKey:Customer;references:CustomerId"`
	Products Product `gorm:"foreignKey:Customer;references:ProductId"`
}

func OneOwnerConstraintTrigger() {
	// Adds trigger constraint that allows to have only one Owner Per Cart, In avoid of merging Orders.
	command := fmt.Sprintf(`CREATE FUNCTION public.check_one_cart_owner() 
	RETURNS TRIGGER 
	LANGUAGE 'plpgsql'
	AS $BODY$ 
	DECLARE updated integer;
	BEGIN UPDATE %s SET owner = owner + 1 WHERE cart.id = NEW.cart_id AND cart.owner < 1;
	GET DIAGNOSTICS addedOwners = ROW_COUNT 
	IF addedOwners = 0 THEN 
	RAISE EXCEPTION 'Cart can have only one owner.'

	END IF;
	RETURN NEW;
	END; 
	$BODY$;
	
	CREATE TRIGGER OwnerCartConstraintTrigger 
	BEFORE INSERT ON public.cart
	FOR EACH ROW EXECUTE PROCEDURE public.check_one_cart_owner();`, "cart")
	Database.Exec(command)
	DebugLogger.Println("Unique Constraint Has Been Integrated.")
}
package main

import (
    "database/sql"
    "fmt"

    _"github.com/lib/pq"
)

const (
  host         = "localhost"
  portOrig     = 5432
  portMig      = 5433
  userOrig     = "old"
  userMig      = "new"
  passwordOrig = "hehehe"
  passwordMig  = "hahaha"
  dbnameOrig   = "old"
  dbnameMig    = "new"
)

type MigratedList struct {
  Accounts []MigratedAccount
}

type MigratedAccount struct {
  Id string
  Name string
  Email string
  FavoriteFlavor string
}

type OriginalList struct {
  Accounts []OriginalAccount
}

type OriginalAccount struct {
  Id string
  Name string
  Email string
}

type ErrorLog struct {
  root *LogEntry
  last *LogEntry
}
type LogEntry struct {
  message string
  next *LogEntry
}

var dbOrig *sql.DB
var dbMig *sql.DB


func main(){
  //Open original database connection
  psqlInfoOrig := fmt.Sprintf("host=%s port=%d user=%s " +
    "password=%s dbname=%s sslmode=disable", host, portOrig, userOrig, passwordOrig, dbnameOrig)
  var dbErr error

  dbOrig,dbErr = sql.Open("postgres", psqlInfoOrig)
  if dbErr != nil {
    panic(dbErr)
  }
  defer dbOrig.Close()

  dbErr = dbOrig.Ping()
  if dbErr != nil{
    panic(dbErr)
  }

  //Open migrated database connection
  psqlInfoMig := fmt.Sprintf("host=%s port=%d user=%s " +
    "password=%s dbname=%s sslmode=disable", host, portMig, userMig, passwordMig, dbnameMig)

  dbMig,dbErr = sql.Open("postgres", psqlInfoMig)
  if dbErr != nil {
    panic(dbErr)
  }
  defer dbMig.Close()

  dbErr = dbMig.Ping()
  if dbErr != nil{
    panic(dbErr)
  }

  var origList OriginalList
  var migList MigratedList

  //Initialize the slice full of the original account structs
  err := origList.initializeOriginal();
  if err != nil {
    panic( err )
  }

  //Initialize the slice full of the migrated account structs
  err = migList.initializeMigrated();
  if err != nil {
    panic( err )
  }

}

/*
Compare Accounts method takes an original account and migrated account struct.
Compares the relevant fields in each struct, and returns an error message
if the fields don't match.
*/
func compareAccounts( orig OriginalAccount, mig MigratedAccount ) ( string, int ) {
  var errorText string = fmt.Sprintf("Conflict between original and migrated version of account %s\n" , orig.Id )
  var errors int = 0

  if orig.Name != mig.Name {
    errors++
    errorText += fmt.Sprintf("Original name: %s\nMigrated name: %s\n\n" , orig.Name , mig.Name )
  }
  if orig.Email != mig.Email {
    errors++
    errorText += fmt.Sprintf("Original email: %s\nMigrated email: %s\n\n" , orig.Email , mig.Email )
  }

  if errors == 0 {
    return "", 0
  }

  return errorText, errors
}

/*
initializeMigrated method is called from the list struct. Queries database for
all accounts and appends them to the slice.
*/
func ( migList *MigratedList ) initializeMigrated() error {
  rows, err := dbMig.Query("SELECT id, name, email, favorite_flavor FROM accounts;")
    if err != nil {
      return err
    }
    defer rows.Close()

    var newAccount MigratedAccount;

    for rows.Next() {
      err = rows.Scan(&newAccount.Id, &newAccount.Name, &newAccount.Email, &newAccount.FavoriteFlavor )
      if err != nil {
        return err
      }

      migList.Accounts = append( migList.Accounts, newAccount )
    }

    return nil
}

/*
initializeOriginal method is called from the list struct. Queries database for
all accounts and appends them to the slice.
*/
func ( origList *OriginalList ) initializeOriginal() error {

  rows, err := dbOrig.Query("SELECT id, name, email FROM accounts;")
    if err != nil {
      return err
    }
    defer rows.Close()

    var newAccount OriginalAccount;

    for rows.Next() {
      err = rows.Scan(&newAccount.Id, &newAccount.Name, &newAccount.Email )
      if err != nil {
        return err
      }

      origList.Accounts = append( origList.Accounts, newAccount )
    }

    return nil
}

/*
getNextOriginal method grabs the first element from the original account list,
deletes the account from the slice, and returns it.
*/
func ( origList *OriginalList ) getNextOriginal() OriginalAccount {
  var nextAccount OriginalAccount

  return nextAccount
}

/*
getNextMigrated method grabs the first element from the original account list,
deletes the account from the slice, and returns it.
*/
func ( migList *MigratedList ) getNextMigrated() error {
  return nil
}

/*
searchMigrated method takes an id string and searches the id values of the
migrated accounts for that account. If it finds the account it deletes the
account from the list and returns it. If it doesn't find the account it returns
an error.
*/
func ( migList *MigratedList ) searchMigrated( id string ) (error, MigratedAccount) {
  var matchingAccount MigratedAccount

  return nil, matchingAccount
}

/*
addLog method takes a string and appends it to the logs linked list.
*/
func ( errorLog *ErrorLog ) addLog( newError string ) {

}

/*
printLog method prints the entire list out in a terminal
*/
func ( errorLog *ErrorLog ) printLog() {

}

/*
saveLog method prints the entire log out into a text document
*/
func ( errorLog *ErrorLog ) saveLog() {

}

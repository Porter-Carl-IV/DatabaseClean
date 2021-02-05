package main

import (
    "database/sql"
    "fmt"
    "errors"
    "os"

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
  Root *LogEntry
  Last *LogEntry
}
type LogEntry struct {
  Message string
  Next *LogEntry
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

  //Initialize logs
  missingLog := ErrorLog{nil, nil}
  corruptedLog := ErrorLog{nil, nil}
  newLog := ErrorLog{nil, nil}

  missingLog.addLog("==========MISSING ACCOUNTS==========\n")
  corruptedLog.addLog("==========CORRUPTED ACCOUNTS==========\n")
  newLog.addLog("==========NEW ACCOUNTS==========\n")

  //Initialize accuracy counts
  missingCount := 0
  corruptedCount := 0
  newCount := 0
  origChecked := 0
  migChecked := 0
  origRecords := len(origList.Accounts)
  migRecords := len(migList.Accounts)

  //Print out total accounts being checked
  fmt.Printf( "Original Accounts: %d\n" , origRecords )
  fmt.Printf( "Migrated Accounts: %d\n" , migRecords )

  //Iterate through original accounts, find missing and corrupted data
  var currentOrigAccount OriginalAccount
  var currentMigAccount MigratedAccount

  for len(origList.Accounts) > 0 {
    currentOrigAccount = origList.getNextOriginal()
    origChecked++

    //Find matching migrated account
    err, currentMigAccount = migList.searchMigrated( currentOrigAccount.Id )
    if err != nil{
      missingCount++;
      missingLog.addLog(
        fmt.Sprintf("Missing record with:\nID: %s\nName:%s\nEmail:%s\n\n", currentOrigAccount.Id , currentOrigAccount.Name , currentOrigAccount.Email ))
    } else {
      //Check for errors with migrated accounts
      errorMessage, errorCount := compareAccounts( currentOrigAccount, currentMigAccount )
      migChecked++
      if errorCount > 0 {
        corruptedCount++
        corruptedLog.addLog( errorMessage )
      }
    }
  }

  //Now that we've removed all matching migrated accounts and all original accounts
  //Iterate through remaining migrated accounts and log them as new accounts
  for len(migList.Accounts) > 0 {
    currentMigAccount = migList.getNextMigrated()

    newCount++
    migChecked++

    newLog.addLog(
      fmt.Sprintf("New record with:\nID: %s\nName:%s\nEmail:%s\nFavorite Flavor:%s\n\n", currentMigAccount.Id , currentMigAccount.Name , currentMigAccount.Email , currentMigAccount.FavoriteFlavor))
  }

  fmt.Printf( "Missing Accounts: %d\n" , missingCount )
  fmt.Printf( "Corrupted Accounts: %d\n" , corruptedCount )
  fmt.Printf( "New Accounts: %d\n" , newCount )

  //Make sure that all records have been iterated over
  if origChecked != origRecords {
    fmt.Printf("Only checked %d out of %d records\n" , origChecked , origRecords)
  }
  if migChecked != migRecords {
    fmt.Printf("Only checked %d out of %d migrated records\n" , migChecked , migRecords)
  }

  //Print logs to file, could instead print to screen with log.printLog instead
  missingLog.saveLog("MissingLog.txt")
  corruptedLog.saveLog("CorruptedLog.txt")
  newLog.saveLog("NewLog.txt")
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

Note: Need to make sure length of list is at least 1 before sending to this method
*/
func ( origList *OriginalList ) getNextOriginal() OriginalAccount {
  var nextAccount OriginalAccount = origList.Accounts[0]

  //Put the last element into the first position, then replace the slice
  //with a version of itself without the last element.
  origList.Accounts[0] = origList.Accounts[ len(origList.Accounts) - 1 ]
  origList.Accounts = origList.Accounts[ :len(origList.Accounts) - 1 ]

  return nextAccount
}

/*
getNextMigrated method grabs the first element from the original account list,
deletes the account from the slice, and returns it.

Note: Need to make sure length of list is at least 1 before sending to this method
*/
func ( migList *MigratedList ) getNextMigrated() MigratedAccount {
  var nextAccount MigratedAccount = migList.Accounts[0]

  //Put the last element into the first position, then replace the slice
  //with a version of itself without the last element.
  migList.Accounts[0] = migList.Accounts[ len(migList.Accounts) - 1 ]
  migList.Accounts = migList.Accounts[ :len(migList.Accounts) - 1 ]

  return nextAccount
}

/*
searchMigrated method takes an id string and searches the id values of the
migrated accounts for that account. If it finds the account it deletes the
account from the list and returns it. If it doesn't find the account it returns
an error.

Note: Need to make sure length of list is at least 1 before sending to this method
*/
func ( migList *MigratedList ) searchMigrated( id string ) (error, MigratedAccount) {
  var matchingAccount MigratedAccount

  for index := range migList.Accounts {
    if migList.Accounts[index].Id == id {
      matchingAccount = migList.Accounts[index]

      //Remove the found account by replacing it with last account, then replacing
      //slice with everything but last element
      migList.Accounts[index] = migList.Accounts[ len(migList.Accounts) - 1 ]
      migList.Accounts = migList.Accounts[ :len(migList.Accounts) - 1 ]

      return nil, matchingAccount
    }
  }

  return errors.New("Account not found"), matchingAccount
}

/*
addLog method takes a string and appends it to the logs linked list.
*/
func ( errorLog *ErrorLog ) addLog( newError string ) {
  newEntry := LogEntry{ newError , nil }

  if errorLog.Root == nil {
    errorLog.Root = &newEntry
    errorLog.Last = errorLog.Root

    return
  }

  errorLog.Last.Next = &newEntry
  errorLog.Last = errorLog.Last.Next
}

/*
printLog method prints the entire list out in a terminal
*/
func ( errorLog *ErrorLog ) printLog() {
  currentLog := errorLog.Root

  for currentLog.Next != nil {
    print( currentLog.Message )
    currentLog = currentLog.Next
  }

  print( currentLog.Message )
}

/*
saveLog method prints the entire log out into a text document
*/
func ( errorLog *ErrorLog ) saveLog( fileName string ) {
  currentLog := errorLog.Root

  file, err := os.Create( fileName )
  if err != nil {
    panic( err )
  }

  defer file.Close()

  for currentLog.Next != nil {
    _, err = file.WriteString( currentLog.Message )

    if err != nil {
      panic( err )
    }
    currentLog = currentLog.Next
  }

  file.WriteString( currentLog.Message )


}

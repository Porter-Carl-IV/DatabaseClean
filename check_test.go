package main

import(
  "testing"
)

func TestCompare( test *testing.T ){
  errorString, errors := compareAccounts( OriginalAccount{ "204" , "Jim" , "jim@yahoo.com" } , MigratedAccount{ "204" , "Jim" , "jim@yahoo.com" , " " } )

  if errors != 0 {
    test.Errorf( "Compare failed, expected no errors got %d errors with message %s\n" , errors , errorString )
  }

  errorString, errors = compareAccounts( OriginalAccount{ "204" , "Jim" , "jim@yahoo.com" } , MigratedAccount{ "204" , "jim" , "jim@yahoo.com" , " " } )

  if errors != 1 {
    test.Errorf( "Compare failed, expected 1 error got %d errors with message %s\n" , errors , errorString )
  }

  errorString, errors = compareAccounts( OriginalAccount{ "204" , "Jim" , "jim@yahoo.com" } , MigratedAccount{ "204" , "jim" , "@yahoo.com" , " " } )

  if errors != 2 {
    test.Errorf( "Compare failed, expected 2 errors got %d errors with message %s\n" , errors , errorString )
  }
}

func TestGetOrig( test *testing.T ){
  testList := OriginalList{ []OriginalAccount{{ "204" , "Jim" , "jim@yahoo.com" } , { "205" , "Jim" , "jim@yahoo.com" } , { "206" , "Jim" , "jim@yahoo.com" } } }

  nextAccount := testList.getNextOriginal()

  if nextAccount.Id != "204" {
    test.Errorf( "Get account failed, expected account 204 got account %s\n" , nextAccount.Id )
  }

  nextAccount = testList.getNextOriginal()

  if nextAccount.Id != "206" {
    test.Errorf( "Get account failed, expected account 206 got account %s\n" , nextAccount.Id )
  }

  nextAccount = testList.getNextOriginal()

  if nextAccount.Id != "205" {
    test.Errorf( "Get account failed, expected account 205 got account %s\n" , nextAccount.Id )
  }

  if( len( testList.Accounts ) != 0 ){
    test.Errorf( "Get account failed, expected 1 record left, got %d\n" , len( testList.Accounts ) )
  }
}

func TestGetMig( test *testing.T ){
  testList := MigratedList{ []MigratedAccount{{ "204" , "Jim" , "jim@yahoo.com" , " " } , { "205" , "Jim" , "jim@yahoo.com" , " " } , { "206" , "Jim" , "jim@yahoo.com" , " " } } }

  nextAccount := testList.getNextMigrated()

  if nextAccount.Id != "204" {
    test.Errorf( "Get account failed, expected account 204 got account %s\n" , nextAccount.Id )
  }

  nextAccount = testList.getNextMigrated()

  if nextAccount.Id != "206" {
    test.Errorf( "Get account failed, expected account 206 got account %s\n" , nextAccount.Id )
  }

  nextAccount = testList.getNextMigrated()

  if nextAccount.Id != "205" {
    test.Errorf( "Get account failed, expected account 205 got account %s\n" , nextAccount.Id )
  }

  if( len( testList.Accounts ) != 0 ){
    test.Errorf( "Get account failed, expected 1 record left, got %d\n" , len( testList.Accounts ) )
  }
}

func TestSearchMig( test *testing.T ){
  testList := MigratedList{ []MigratedAccount{{ "204" , "Jim" , "jim@yahoo.com" , " " } , { "205" , "Jim" , "jim@yahoo.com" , " " } , { "206" , "Jim" , "jim@yahoo.com" , " " } } }

  err, foundAccount := testList.searchMigrated("204")

  if foundAccount.Id != "204" {
    test.Errorf( "Search account failed, expected account 204 account %s with error %s\n" , foundAccount.Id , err )
  }

  err, foundAccount = testList.searchMigrated("205")

  if foundAccount.Id != "205" {
    test.Errorf( "Search account failed, expected account 205 account %s with error %s\n" , foundAccount.Id , err  )
  }

  err, foundAccount = testList.searchMigrated("206")

  if foundAccount.Id != "206" {
    test.Errorf( "Search account failed, expected account 206 got account %s with error %s\n" , foundAccount.Id , err  )
  }

  if( len( testList.Accounts ) != 0 ){
    test.Errorf( "Search account failed, expected 1 record left, got %d\n" , len( testList.Accounts ) )
  }
}

func TestAddLog( test *testing.T ) {
  testLog := ErrorLog{ nil , nil }

  testLog.addLog( "First Log" )

  if testLog.Root.Message != "First Log" {
    test.Errorf( "Add log action failed, expected \"First Log\" but got \"%s\"\n" , testLog.Root.Message )
  }

  testLog.addLog( "Second Log" )

  if testLog.Root.Message != "First Log" {
    test.Errorf( "Second Add log action failed, expected \"First Log\" still in root but got \"%s\"\n" , testLog.Root.Message )
  }

  if testLog.Last.Message != "Second Log" {
    test.Errorf( "Second Add log action failed, expected \"Second Log\" in last but got \"%s\"\n" , testLog.Last.Message )
  }

  testLog.addLog( "Third Log" )

  if testLog.Root.Message != "First Log" {
    test.Errorf( "Third Add log action failed, expected \"First Log\" still in root but got \"%s\"\n" , testLog.Root.Message )
  }

  if testLog.Root.Next.Message != "Second Log" {
    test.Errorf( "Third Add log action failed, expected \"Second Log\" next of root but got \"%s\"\n" , testLog.Root.Next.Message )
  }

  if testLog.Last.Message != "Third Log" {
    test.Errorf( "Third Add log action failed, expected \"Third Log\" in last but got \"%s\"\n" , testLog.Last.Message)
  }
}

/*
This method isn't really a conventional test. To make it one I would need to compare against the stdout bitstream.
However moving forward with the assumption that in this application it's not worth taking the extra time to implement
the test that way.
*/
func TestPrintLog( test *testing.T ) {
  testLog := ErrorLog{ nil , nil }

  testLog.addLog( "First Log\n" )
  testLog.addLog( "Second Log\n" )
  testLog.addLog( "Third Log\n" )

  print("Expect:\nFirst Log\nSecond Log\nThird Log\nGot:\n")

  testLog.printLog()
}

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

  if( len( testList.Accounts ) != 1 ){
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

  if( len( testList.Accounts ) != 1 ){
    test.Errorf( "Get account failed, expected 1 record left, got %d\n" , len( testList.Accounts ) )
  }
}

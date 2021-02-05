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

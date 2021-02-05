# DatabaseClean
Golang tool that takes two public databases, compares for differences, and outputs results.


To run code:
publish the docker container of both databases to your own local host on ports 5432 (original db) and 5433 (migrated db) using commands:

docker run -p 5433:5432 guaranteedrate/homework-post-migration:1607545060-a7085621

docker run -p 5432:5432 guaranteedrate/homework-pre-migration:1607545060-a7085621

then migrate to directory and use command

go run check.go

To test code:
migrate to directory and use command

go test -v

Assumptions:
-Because this is a progam that is only run a couple times to find problems, efficiency is less important than clarity of execution.
Therefore I've focused writing the code in a way that is readable and modular so individual pieces can be readily understood and tested.

-The person running this code has access to the databases in question, and therefore knows the number of records in each.
Therefore instead of writing a unit test to make sure that every record on the database is tested, I printed out the number of records tested
for the user to check.

-The person running this code is a developer, and so they don't need command line flags to change the output mode. They can instead choose whether to
call the printLog, saveLog, or both functions at the end.

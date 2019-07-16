# cifrest

Originally this was a separate golang application used to expose the CIF data as rest endpoints.

However with the advent of the dbrest utility within our uktransport repository all of the golang code is now obsolete
as that one utility does everything we want. All that remains is a set of postgresql functions and the configuration
file to setup the endpoints for that utility.

## Installation

To install you need to read in each of the sql scripts in this directory into the same database & schema that was created by
cifimport. This will add all of the necessary postgresql functions needed for each endpoint.

You then need to copy config-example.yaml from this directory and edit the datasource settings to match your database.

You can then run the dbrest utility pointing it to that new yaml file and your CIF database is now exposed as rest endpoints.


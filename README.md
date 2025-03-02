# Assignment 1

Assignment 1 for Prog2005 Cloud - Tharald Roland Sørensen

## Guidance
The service has 3 endpoints:

- /countryinfo/v1/info
- /countryinfo/v1/population
- /countryinfo/v1/status

## General info about endpoitns:

../info is to get general information about a given country.

../population is to get information about a countries' population.

../status is to get the general status of the used API's, the services uptime and version.

### Note

../info and ../population both takes a /{country_code} parameter at the end, this is the Iso2 code for that country: https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2

../info and ../population both have a query parameter:
### info:

?limit=(0->10)      e.g., limit=7

Without a given limit, the service will display all cities.

### population:

?limit=xxxx-xxxx        e.g., limit=2010-2015

Without a given limit, the service will display all years

### Note

limit only takes 2 parameters divided by: "-", no more, no less.

##
# [blackblade-batch](https://blackblade.ca)

The batch for fetching card data for [blackblade](https://github.com/BrandonWade/blackblade).

## Getting Started

### Prerequisites

1. Clone the [blackblade-infrastructure](https://github.com/BrandonWade/blackblade-infrastructure) repo
2. Run `docker-compose up --build -d` to start the db container

### Fetching Card Data

1. Run `docker-compose up --build` to start the batch container
2. [Optional] The batch will take several minutes to run. Once it completes, confirm the card data is visible in the database

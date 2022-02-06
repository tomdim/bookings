# Bookings and Reservations

This is the repository for my bookings and reservations project.

### Built With

* [Golang 1.17](https://go.dev/)
  * Uses the [chi router](https://github.com/go-chi/chi/v5)
  * Uses [Alex Edwards SCS](https://github.com/alexedwards/scs/v2) session management
  * Uses [nosurf](https://github.com/justinas/nosurf)
  * Uses [govalidator](https://github.com/asaskevich/govalidator) Package of validators and sanitizers for strings, numerics, slices and structs
* [Bootstrap](https://getbootstrap.com)
* [Javascript](https://www.javascript.com/)
  * Uses [notie](https://jaredreich.com/notie/)
  * Uses [sweetalert](https://sweetalert2.github.io/)
  * Uses [vanillajs-datepicker](https://mymth.github.io/vanillajs-datepicker)


## Getting Started

Instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps:

### Prerequisites

Have Golang installed into your system. (Documentation: https://go.dev/doc/install)

### Installation - Running the app locally

1. Clone the repo
   ```shell
   git clone https://github.com/tomdim/bookings.git
   ```
2. Go to project directory
   ```shell
   cd bookings
   ```
3. Download dependencies 
   ```shell
   go mod download
   ```
4. Run server
   ```shell
   make run
   ```
5. Open your browser and go to http://localhost:8889/.

### Running the app with Docker
Steps 1 & 2 apply here as well.

Instead of steps 3 and 4, you will run:
```shell
make build
```
Then, open your browser and go to http://localhost:8889/.

### Run tests locally
In order to run tests locally, go to project directory and run:
```sh
make test
```

### Makefile Documentation
In order to see the make commands you can run, go to terminal and run:
```shell
make help
```
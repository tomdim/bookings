# Bookings and Reservations

This is the repository for my bookings and reservations project.

### Built With

* [Golang 1.17](https://go.dev/)
  * Uses the [chi router](https://github.com/go-chi/chi/v5)
  * Uses [Alex Edwards SCS](https://github.com/alexedwards/scs/v2) session management
  * Uses [nosurf](https://github.com/justinas/nosurf)
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

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/tomdim/bookings.git
   ```
2. Go to project directory
   ```sh
   cd bookings
   ```
3. Run server
   ```sh
   go run cmd/web/*.go
   ```
4. Open your browser and go to http://localhost:8889/

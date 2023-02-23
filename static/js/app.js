function Prompt() {
  let toast = (c) => {
    const { msg = "", icon = "success", position = "top-end" } = c;

    const Toast = Swal.mixin({
      toast: true,
      position: "top-end",
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
      didOpen: (toast) => {
        toast.addEventListener("mouseenter", Swal.stopTimer);
        toast.addEventListener("mouseleave", Swal.resumeTimer);
      },
    });

    Toast.fire({
      icon: icon,
      title: msg,
    });
  };

  let success = (c) => {
    const { msg = "", title = "", footer = "" } = c;

    const Toast = Swal.mixin({
      toast: true,
      position: "top-end",
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
      didOpen: (toast) => {
        toast.addEventListener("mouseenter", Swal.stopTimer);
        toast.addEventListener("mouseleave", Swal.resumeTimer);
      },
    });

    Toast.fire({
      icon: "success",
      title: title,
      text: msg,
    });
  };

  let error = (c) => {
    const { msg = "", title = "", footer = "" } = c;

    const Toast = Swal.mixin({
      toast: true,
      position: "top-end",
      showConfirmButton: false,
      timer: 3000,
      timerProgressBar: true,
      didOpen: (toast) => {
        toast.addEventListener("mouseenter", Swal.stopTimer);
        toast.addEventListener("mouseleave", Swal.resumeTimer);
      },
    });

    Toast.fire({
      icon: "error",
      title: title,
      text: msg,
    });
  };

  async function custom(c) {
    const { icon = "", msg = "", title = "", showConfirmButton = true } = c;
    const { value: result } = await Swal.fire({
      icon: icon,
      title: title,
      html: msg,
      focusConfirm: false,
      showCancelButton: true,
      showConfirmButton: showConfirmButton,
      // willOpen: () => {
      //   const elem = document.getElementById("reservation-dates-modal");
      //   const rp = new DateRangePicker(elem, {
      //     format: "yyyy-mm-dd",
      //     showOnFocus: true,
      //   });
      // },
      preConfirm: () => {
        return [
          document.getElementById("start").value,
          document.getElementById("end").value,
        ];
      },
      didOpen: () => {
        if (c.didOpen != undefined) {
          c.didOpen();
        }
      },
    });

    if (result) {
      if (result.dismiss !== Swal.DismissReason.cancel) {
        if (result.value !== "") {
          if (c.callback !== undefined) {
            c.callback(result);
          }
        } else {
          c.callback(false);
        }
      } else {
        c.callback(false);
      }
    }
  }

  return {
    toast: toast,
    success: success,
    error: error,
    custom: custom,
  };
}

// Function to display the prompt for reservation
const ShowDatesBox = (id) => {
  document
    .getElementById("check-availability-button")
    .addEventListener("click", () => {
      let html = `
        <form id="check-availability" action = "" method = "post" novalidate class = "needs-validation">
          <div class = :form-row">
            <div class =  "col">
              <div class = "form-row" id =  "reservation-dates-modal">
                <div class = "col">
                  <input disabled required class = "form-control" type="date" name="start" id="start" placeholder="Arrival">
                </div>
                <div class = "col">
                  <input disabled required class = "form-control" type="date" name="end" id="end" placeholder="Departure">
                </div>
              </div>
            </div>
          </div>
        </form>
        `;
      attention.custom({
        title: "Enter your dates here!",
        msg: html,
        didOpen: () => {
          document.getElementById("start").removeAttribute("disabled");
          document.getElementById("end").removeAttribute("disabled");
        },
        callback: function (result) {
          let form = document.getElementById("check-availability");
          let formData = new FormData(form);

          formData.append("csrf_token", "{{.CSRFToken}}");
          formData.append("room_id", id);

          fetch("/search-availability-json", {
            method: "post",
            body: formData,
          })
            .then((response) => response.json())
            .then((data) => {
              if (data.ok) {
                link =
                  "/book-room?id=" +
                  data.room_id +
                  "&s=" +
                  data.start_date +
                  "&e=" +
                  data.end_date;

                attention.custom({
                  icon: "success",
                  msg:
                    "<p>Room is Available </p>" +
                    '<p><a href="' +
                    link +
                    '"' +
                    'class="btn btn-primary">Book Now!</a></p>',
                  showConfirmButton: false,
                });
              } else {
                attention.error({
                  msg: "No availibility!",
                });
              }
            });
        },
      });
    });
};

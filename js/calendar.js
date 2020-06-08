$.ajax({
  url: `/_getDateEvent`,
  contentType: "application/json;charset=UTF-8",
  type: "POST",
  xhrFields: {
    withCredentials: true,
  },
}).done(function (data) {
  var events = $.parseJSON(data);
  // Parse our events into an object called events that will later be used to initialize FullCalendar
  /* var events = []; */
  // console.log(events);

  //   $.each(data, function (i, v) {
  //     events.push(v);
  //   });

  var calendarEl = document.getElementById("calendar");

  var calendar = new FullCalendar.Calendar(calendarEl, {
    plugins: ["interaction", "dayGrid", "timeGrid", "list"],
    header: {
      left: "prev,next today",
      center: "title",
      right: "dayGridMonth,listMonth",
    },
    eventLimit: true,
    navLinks: true, // can click day/week names to navigate views
    businessHours: true, // display business hours
    editable: false,
    events: events,
  });

  calendar.render();
});

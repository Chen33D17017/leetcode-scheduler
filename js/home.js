$(document).ready(function () {
  $.ajax({
    url: `/_getUndo`,
    contentType: "application/json;charset=UTF-8",
    type: "POST",
    xhrFields: {
      withCredentials: true,
    },
    success: function (result) {
      objs = $.parseJSON(result);
      for (obj of objs) {
        addCard(obj);
        console.log(obj);
      }
      rebund();
    },
    error: function (err) {
      $("#add-msg").text(err);
    },
  });
  // rebund();
  $("#delete-problem").on("click", function () {
    var target = $(this).attr("target");
    var deleteDom = `#${$(this).attr("problem")}`;
    $.ajax({
      url: `/deleteLog/${target}`,
      type: "DELETE",
      xhrFields: {
        withCredentials: true,
      },
      success: function (result) {
        $(deleteDom).parent().remove();
        $("#close-delete").click();
      },
      error: function (err) {
        console.log("delete error");
      },
    });
  });

  $("#stop").on("click", function () {
    let stopwatch = $(this).parent().parent().children("h1");
    let minutes = stopwatch.children("#minutes").text();
    let seconds = stopwatch.children("#seconds").text();
    let hundredths = stopwatch.children("#hundredths").text();
    let time = `${minutes}:${seconds}:${hundredths}`;
    var target = $("#problem-id").text();
    $(target).text(time);
  });

  $(".type-btn").on("click", function () {
    $(this).parent().children("#problem-type").toggle();
  });

  $("#add-problem").on("click", function () {
    var addContent = $("#add-content").val();
    var msg = addContent.length > 0 ? addContent : "Enter content";
    $("#add-msg").text(msg);
    $.ajax({
      url: `/addProblem/${msg}`,
      contentType: "application/json;charset=UTF-8",
      type: "POST",
      xhrFields: {
        withCredentials: true,
      },
      success: function (result) {
        obj = $.parseJSON(result);
        addCard(obj);
        rebundPlay();
      },
      error: function (err) {
        $("#add-msg").text(err);
      },
    });
  });
});

function rebund() {
  $(".play-button").on("click", function () {
    $("#start").click();
    let problemName = $(this)
      .parent()
      .parent()
      .children("a")
      .children(".problem-name")
      .text();
    $("#target-problem").text(problemName);
    $("#problem-id").text("#" + $(this).attr("target"));
  });

  $(".close-p").on("click", function () {
    var closeTarget = $(this).attr("close-target");
    var closeProblem = $(this).attr("close-problem");
    // console.log(closeTarget);
    console.log(closeProblem);
    $("#delete-problem").attr("target", closeTarget);
    $("#delete-problem").attr("problem", closeProblem);
  });

  $(".feedback").on("click", function () {
    data = {};
    var elem = $(this);
    var target = $(this).attr("target");
    console.log(target);
    data["level"] = $(this).text().toLowerCase();
    data["costTime"] = $(this).parent().prev("div").children("span").text();
    console.log(data);
    $.ajax({
      url: `/DoneProblem/${target}`,
      contentType: "application/json;charset=UTF-8",
      type: "POST",
      xhrFields: {
        withCredentials: true,
      },
      data: JSON.stringify(data),
      success: function (result) {
        elem.parent().parent().parent().parent().parent().parent().remove();
        console.log("delete this block");
      },
      error: function (err) {
        $("#add-msg").text("Bad Try");
      },
    });
  });
}

function addCard(obj) {
  $("#accordion").prepend(
    `
      <div class="card">
          <div class="card-header ${obj["Difficulty"].toLowerCase()}" id="p-${
      obj["id"]
    }">
            <h5 class="mb-0">
              <button class="header-font btn btn-link" data-toggle="collapse" data-target="#collapse-${
                obj["id"]
              }" aria-expanded="true" aria-controls="collapse-${obj["id"]}">
                ${obj["id"]}. ${obj["problem_name"]}
              </button>
              <button type="button" class="close close-p" data-toggle="modal" data-target="#deleteModal" close-target="${
                obj["LogID"]
              }" close-problem="p-${obj["id"]}">
                <span aria-hidden="true">&times;</span>
              </button>
            </h5>
          </div>
          <div id="collapse-${
            obj["id"]
          }" class="collapse " aria-labelledby="headingOne" data-parent="#accordion">
            <div class="card-body">
              <div class="clearfix">
                <a href="${obj["url"]}">
                  <span class="problem-name">${obj["problem_name"]}</span>
                </a>
                <span class="float-right">
                  <button class="play-button" target="ft-${obj["id"]}">
                    <img class="icon" src="../static/img/play.png" alt="play.png">
                  </button>
                </span>
              </div>
              <div class="row" style="float: clear;">
                <div class="col-4">
                  <p><b>Difficulty:</b> ${obj["Difficulty"]}</p>
                </div>
                <div class="col-4">
                  <p><b>Review Level:</b> ${obj["ReviewLevel"]}</p>
                </div>
                <div class="col-4">
                  <p><b>Deadline:</b> ${obj["Deadline"]}</p>
                </div>
                <div class="col-6">
                  <button class="type-btn">
                    <div style="font-weight: bold; ">Problem Type: </div>
                  </button>                      
                  <div id="pt-${obj["id"]}" style="display: none;">
                    <button id="add-category">+</button>
                  </div>
                </div>
                <div class="col-6">
                  <div style="font-weight: bold; "> Finish Time: <span id="ft-${
                    obj["id"]
                  }"> 00:00:00 </span></div>
                  <div>
                    <button class="feedback easy" target="${
                      obj["LogID"]
                    }">Easy</button> 
                    <button class="feedback medium" target="${
                      obj["LogID"]
                    }">Medium</button>
                    <button class="feedback hard" target="${
                      obj["LogID"]
                    }">Hard</button>
                  </div>
                </div>
            </div>
          </div>
        </div>
      `
  );
}

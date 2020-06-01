$(document).ready(function () {
  $(".input").focus(function () {
    $(this).parent().parent().addClass("focus");
  });

  $(".input").blur(function () {
    if (!$(this).val()) {
      $(this).parent().parent().removeClass("focus");
    }
  });

  // check whether email valid
  $("#r-email").on("change", function () {
    var email = $(this).val();
    var emailErr = $("#email-err > h5");
    $.ajax({
      url: "/_checkUserExist",
      type: "post",
      data: {
        msg: email,
      },
      success: function (data) {
        if (data === "OK") {
          $("#regist").removeAttr("disabled");
          emailErr.text("");
        } else {
          emailErr.text("Email already been used");
          $("#regist").attr("disabled", "disabled");
        }
      },
    });
    emailErr.text("");
  });

  $("#password-confirm").on("change", function () {
    let password = $("#password").val();
    let psErr = $("#password-err > h5");
    let pscErr = $("#password-confirm-err > h5");
    if (password != $(this).val()) {
      pscErr.text("Password Unmatch, please try again");
      psErr.text("Password Unmatch, please try again");
      $("#regist").attr("disabled", "disabled");
    } else {
      pscErr.text("");
      psErr.text("");
      $("#regist").removeAttr("disabled");
    }
  });
});

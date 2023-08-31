const myFunction = () => {
    mins = new Date().getMinutes();
    if (mins < 10) {
        mins = "0" + mins;
    }
    var time = new Date().getHours() + ":" + mins;
    $("#dash-time-txt").text(time);

};

$(".goto-chat").click(function(){
        $(".loading").fadeIn();
    setTimeout(function(){
        history.pushState({}, null, "/chat")
        document.title = "Probe / Chat";

        document.querySelector("body").innerHTML='<object class="embed" type="text/html" data="/chat" ></object>';
        setTimeout(function(){
                $(".loading").fadeOut();
        }, 3000)
    }, 2000)
})
$(".goto-map").click(function(){
    $(".loading").fadeIn();
setTimeout(function(){
    history.pushState({}, null, "/map")
    document.title = "Probe / Map";

    document.querySelector("body").innerHTML='<object class="embed" type="text/html" data="/map" ></object>';
    setTimeout(function(){
            $(".loading").fadeOut();
    }, 3000)
}, 2000)
})
$(".goto-home").click(function(){
    $(".loading").fadeIn();
setTimeout(function(){
    history.pushState({}, null, "/")
    document.title = "Probe / Home";
    document.querySelector("body").innerHTML='<object class="embed" type="text/html" data="/" ></object>';
    setTimeout(function(){
            $(".loading").fadeOut();
    }, 3000)
}, 2000)
})

$(".goto-resources").click(function(){
    $(".loading").fadeIn();
setTimeout(function(){
    history.pushState({}, null, "/resources")
    document.title = "Probe / Resources";
    document.querySelector("body").innerHTML='<object class="embed" type="text/html" data="/resources" ></object>';
    setTimeout(function(){
            $(".loading").fadeOut();
    }, 3000)
}, 2000)
})


$(".goto-request").click(function(){
    $(".loading").fadeIn();
setTimeout(function(){
    history.pushState({}, null, "/request")
    document.title = "Probe / Request";
    document.querySelector("body").innerHTML='<object class="embed" type="text/html" data="/request" ></object>';
    setTimeout(function(){
            $(".loading").fadeOut();
    }, 3000)
}, 2000)
})
function ijwei(){
    $(".login-arow").animate({
        width: "100%",
        backgroundColor: "#c4c4c4",
        left: "0",
        margin: "0 0 0 0",
        height: "4em"
    }, 1000)
}
$(".login-submit").click(function(){
    event.preventDefault()
    
    $.ajax({
        url: "/login/verify",
        type: "POST",
        data: {
            username: parseInt($("#username").val()),
            password: $("#password").val()
        },
        success: function(response) {
            ijwei()
            setTimeout(function(){
                window.location.href = "/"
            }, 1500)
            console.log(response)
        },
        error: function(error) {
            if (error.status === 409){
                $(".success-cipher").fadeIn(500);
                $("#sc-title").text("Agent cipher different from ID")
                setTimeout(function(){
                    $(".success-cipher").fadeOut();
                }, 4000)
            }
            if (error.status === 417){
                $(".success-cipher").fadeIn(500);
                $("#sc-title").text("Incorrect password")
                setTimeout(function(){
                    $(".success-cipher").fadeOut();
                }, 4000)
            }else{
                console.log(error)
                
            }
            console.log(error)
                // $(".success-cipher").fadeIn(500);
                // $("#sc-title").text("Incorrect information")
                // setTimeout(function(){
                //     $(".success-cipher").fadeOut();
                // }, 4000)
        }
})
})

$(".login-form").submit(function(){
    event.preventDefault()
    e.preventDefault()
})



$(".cipher-submit").click(function(){
    event.preventDefault()
    var fileInput = $("#upload")[0];
    if (fileInput.files.length > 0) {
        var file = fileInput.files[0];
        if (file.name.endsWith(".txt")) {
            var formData = new FormData();
            formData.append("file", file);

            $.ajax({
                url: "/check/cipher",
                type: "POST",
                data: formData,
                processData: false,
                cache: false,
                contentType: false,
                success: function(response) {
                    $(".cipher-upload").fadeOut();
                    console.log(response)
                    setTimeout(function(){
                        $(".success-cipher").fadeIn(500);
                        setTimeout(function(){
                            $(".success-cipher").fadeOut();
                        }, 4000)
                    })
                },
                error: function(error) {
                    if (error.status === 409){
                        console.log(error)
                        setTimeout(function(){
                            $(".success-cipher").fadeIn(500);
                            $("#sc-title").text("Incorrect cipher")
                            setTimeout(function(){
                                $(".success-cipher").fadeOut();
                            }, 4000)
                        })
                    } 
                    else{
                        console.log(error)
                    }
                }
            });
        } else {
            console.log("Please select a .txt file");
        }
    }
})

$(".mission-viewmore").click(function(){
    $(".missions-fullscreen").fadeIn();
})

$(".fbg").click(function(){
    $(".missions-fullscreen").fadeOut()
})

$(".rkkclose").click(function(){
    $(".req-missions").fadeOut()
})

$(".qwiejqwioee").click(function(){
    $(".req-missions").fadeIn()
})
setInterval(myFunction, 100); 

$('.req-mission').click(function(){
    const missionId = this.id.replace('req-mission-', ''); 
    $.ajax({
        url: `/request/mission/${missionId}`,
        type: "POST",
        success: function(response) {
            setTimeout(function(){
                $(".success-cipher").fadeIn(500);
                $("#sc-title").text(`Added mission`)
                setTimeout(function(){
                    $(".success-cipher").fadeOut();
                }, 4000)
                console.log(response)
            })
        },
        error: function(error) {
                setTimeout(function(){
                    $(".success-cipher").fadeIn(500);
                    $("#sc-title").text(`error ${error.responseJSON.error}`)
                    setTimeout(function(){
                        $(".success-cipher").fadeOut();
                    }, 4000)
                })
                console.log(error)
        }
    });
});
var reqAmt = 4;
$("#req-plus").click(function(){
    if (reqAmt > 9){
        reqAmt =9;
    }
    reqAmt++;
    $("#req-amt").text(reqAmt)
})
$("#req-minus").click(function(){
    if (reqAmt < 2){
        reqAmt = 2;
    }
    reqAmt--;
    $("#req-amt").text(reqAmt)
})
var reqUrgency
$(".req-step-check").click(function(){
    reqUrgency = (this.id)
    var $this = $(this);
    $(".req-step-check").animate({
        backgroundColor: "#1e1e1e"
    }, 300)
    setTimeout(function(){

        $(".req-step-check").removeClass("checked");
    }, 400)
    $this.animate({
        backgroundColor: "#c1c1c1"
    }, 300)
    setTimeout(function(){
        $this.addClass("checked")
    }, 400)
})
console.log(reqUrgency)

function systumPeSystum(){
    if ("geolocation" in navigator) {
        navigator.geolocation.getCurrentPosition(
          (position) => {
            const lat = position.coords.latitude;
            const lng = position.coords.longitude;
      
            console.log(`Latitude: ${lat}, longitude: ${lng}`);
            if (confirm("Activate mayday? Your access will be revoked to prevent compromise")){
                $.ajax({
                    url: "/notification/mayday",
                    type: "POST",
                    data: {
                        lat: lat,
                        lng: lng
                    },
                    success: function(response) {
                        console.log(response)
                        window.location.href = "/"
                    },
                    error: function(error) {
                        console.log(error)
                    }
                })
            }
          },
          (error) => {
            console.error("Error getting user location:", error);
          }
        );
      } else {
        console.error("Geolocation is not supported by this browser.");
      }
}
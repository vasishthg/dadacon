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

const missionCards = document.querySelectorAll('.dash-mission');
const prevButton = document.querySelector('#m-prev-btn');
const nextButton = document.querySelector('#m-next-btn');

let currentMissionIndex = 0;

function showMission(index) {
    missionCards.forEach((card, i) => {
        card.style.display = i === index ? 'block' : 'none';
    });
}

function prevMission() {
    currentMissionIndex = (currentMissionIndex - 1 + missionCards.length) % missionCards.length;
    showMission(currentMissionIndex);
}

function nextMission() {
    currentMissionIndex = (currentMissionIndex + 1) % missionCards.length;
    showMission(currentMissionIndex);
}

prevButton.addEventListener('click', prevMission);
nextButton.addEventListener('click', nextMission);

// Show the first mission card initially
showMission(currentMissionIndex);
$(`req-mission-{{id}}`).click(function(){
    const missionId = this.id.replace('req-mission-', ''); // Extract the ID from the element ID
    $.ajax({
        url: `/request/mission/${missionID}`,
        type: "POST",
        success: function(response) {
            console.log(response)
        },
        error: function(error) {
            console.log(error)
        }
    })
})
const myFunction = () => {
    var time = new Date().getHours() + ":" + new Date().getMinutes();
    console.log(time);
    $("#dash-time-txt").text(time);
};

setInterval(myFunction, 100); 
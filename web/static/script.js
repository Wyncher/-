function downloadBlob(byteArr, name = 'file.txt') {
    //base64 to binary
    const decodedData = atob(byteArr);
    //binary to array binary
    const arrayBuffer = new Uint8Array(decodedData.length);
    for (let i = 0; i < decodedData.length; i++) {
        //decode everyone char
        arrayBuffer[i] = decodedData.charCodeAt(i);
    }
    //make blob
    const blob = new Blob([arrayBuffer]);
    //make blob url(content of file)
    const blobUrl = URL.createObjectURL(blob);
    //make link in DOM
    const link = document.createElement("a");

    link.href = blobUrl;
    link.download = name;
    document.body.appendChild(link);
    link.dispatchEvent(
        //make event
        new MouseEvent('click', {
            bubbles: true,
            cancelable: true,
            view: window
        })
    );

    document.body.removeChild(link);
}
function refresh(){
    $.ajax({
        url: "http://localhost:1323/refresh",
        method: "POST",
        data : {},
        success : function(data) {
            if (data !=null){
                //If container contains light messages then new message must be in light
                let theme
            if ($('.messages').find('.messageRight').length >0) {
                theme = "";
            }else{
                theme = "dark/";
                }
            //for every message
                for (let i = 0; i < data.length; i++) {
                    if (data[i]["Pos"] == "r") {
                        //add LIGHT/DARK outgoing message to html
                        $(' <div class="messageRight'+theme.slice(0,4)+'">' + data[i]["Text"] +'<span>' + data[i]["DateStr"] +
                            '</span><span style="cursor:pointer;" onclick="downloadBlob(' + data[i]["File"] +
                            ', \'' + data[i]["FileName"] + '\');">' + data[i]["FileName"] +
                            '</span><div class="messageBottom"><a href="/chat?delete=' + data[i]["ID"] +
                            '"><img src="static/'+theme+'trash.png" alt=""></a><a href="/chat?favourite=' + data[i]["ID"] +
                            '"><img src="static/'+theme+'star.png" alt=""></a></div></div>').appendTo('.messages');

                    }else{
                        //add LIGHT/DARK incoming message to html
                        $(' <div class="messageLeft'+theme.slice(0,4)+'">'+data[i]["Text"] + '<span>' + data[i]["DateStr"] +
                            '</span><span style="cursor:pointer;" onclick="downloadBlob(' + data[i]["File"] +
                            ', \'' + data[i]["FileName"] + '\');">' + data[i]["FileName"] +
                            '</span><div class="messageBottom"><a href="/chat?delete=' + data[i]["ID"] +
                            '"><img src="static/'+theme+'trash.png" alt=""></a><a href="/chat?favourite=' + data[i]["ID"] +
                            '"><img src="static/'+theme+'star.png" alt=""></a></div></div>').appendTo('.messages');
                    }
                }
            }
            },
    });
} function loadNewMessage(){
    $.ajax({
        url: "http://localhost:1323/loadnewmessage",
        method: "POST",
        data : {},
        success : function(data) {

            if (data != null) {
                for (let i = 0; i < data.length; i++) {
                    const messageCount = document.createElement('h2');
                    messageCount.className = `newMessage`;
                    messageCount.innerHTML = data[i].MCounter;
                    usercard = $('#' + data[i].UserID).parents('a');
                    usercard.children().not('h1').not('img').remove();
                    usercard[0].appendChild(messageCount)
                }
            }
        },
    });
}

function fillUsers(data){
    if (data !=null) {
        $('#userListContainer').html('<div class="userCard"><a href="/chat?favouritechats=1"><img alt="Красивый логотип" id="" alt="" src="static/star.png"><h1>Избранные</h1></a></div>');
        for (let i = 0; i < data.length; i++) {
            $('<div class="userCard">' + '<a href="?recipient=' + data[i].Id + '"><img alt="Красивый логотип"id="'+data[i].Id+'" alt=""><h1>' + data[i].Username + '</h1></a></div>').appendTo('#userListContainer');
            $('#' +data[i].Id).attr('src', `data:image/jpeg;base64,`+data[i].Logo);
        }
    }
}
function loadUserConnects(){
    $.ajax({
        url: "http://localhost:1323/loaduserconnects",
        method: "POST",
        data : {},
        success : function(data) {
            fillUsers(data)
        },
    });
}

//search panel

$( document ).ready(function() {
    //scroll down
    $('.messages').animate({
        scrollTop: $('.messages')[0].scrollHeight

    }, 0);

    //hide smile bar
    $('.smileSlideOuter').hide()
    //add smile
    $("div.smileSlide a").on('click',function () {
        var inputEle = document.querySelectorAll('#text-input');
        inputEle[0].value += $(this).text();
        console.log(inputEle[0].value)
    });
    //open/close smile bar
    $('.smileBtn').click(function(){
        $('.smileSlideOuter').slideToggle(500);
    });
    //refresh messages
    setInterval('refresh()',5000);
    loadUserConnects()
    setInterval('  loadNewMessage()',10000);


    $('#send').click(function() {
        $.ajax({
            url: "http://localhost:1323/chat",
            method: "POST",
            data : { search: $('#searchInput').val()},
            success : function(data) {
                fillUsers(data)
            },
        });
    });
});
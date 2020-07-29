var name = document.cookie.replace(/(?:(?:^|.*;\s*)name\s*\=\s*([^;]*).*$)|^.*$/, "$1");
$("#name").html(name);

var page_number=0; 

var socket = new WebSocket("ws://localhost:4200/socket"); 

socket.onopen = function () {
	console.log("Status: Connected");
};

socket.onmessage = function (e) {
	var obj = JSON.parse(e.data);
	load_com(obj.Title, obj.Description);
	//console.log("get msg - "+e.data);	
};

    function send() {
        //socket.send(input.value);
        //input.value = "";
	}

function get_hedings(){
	$("#reason").html("");
	$("#inner").html("");
	document.getElementById("titles").style.display = "none"; 
	$.ajax({       
			url: 'http://localhost:4200/titles',	
			type: "GET",
			dataType: "json",
			success: function(data){
				for(i = 0; i < data.length; i++){
					$("#reason").append("<p><a href=#/page"+String(parseInt(parseInt(i)+1))+" onclick=\"get_page("+parseInt(i+1)+")\">"+String(i+1)+". "+data[i].Title+"</a></p>");
				}
			},
			failure: function(errMsg) {
				alert(errMsg);
			}
		});
}

function get_page(_page){
	//console.log("_page "+_page);
	page_number = _page;
	if (_page == 0){
		document.getElementsByClassName("img_left")[0].style.display = "none";
		document.getElementsByClassName("img_right")[0].style.display = "none"; 		
		get_hedings();
		$("#box").empty();
	}
	else {
		if(_page != 1){
			document.getElementsByClassName("img_left")[0].style.display = "block";
		}
		else{
			document.getElementsByClassName("img_left")[0].style.display = "none";
		}
		if(_page != 3){
			document.getElementsByClassName("img_right")[0].style.display = "block"; 
		}
		else{
			document.getElementsByClassName("img_right")[0].style.display = "none"; 
		}
		document.getElementById("titles").style.display = "block";
		get_content(_page);
		sleep(100).then(() => {get_comments(_page);});
	}
}

function get_content(_page){
	$.ajax({       
			url: 'http://localhost:4200/content/'+String(_page),	
			type: "GET",
			dataType: "json",
			success: function(data){
				load_content(data.Title, data.Content);
			},
			failure: function(errMsg) {
				alert(errMsg);
			}
		});
}

function load_content(name, _text){
	$("#reason").html("<b>"+name+"</b>");
	$("#inner").html(_text);
}

function get_comments(_page){
	$.ajax({       
			url: 'http://localhost:4200/comms',	
			type: "POST",
			data: JSON.stringify({Pag: String(_page)}),
			dataType: "json",
			success: function(data){
				for(i = 0; i < data.length; i++){
					load_com(data[i].Title, data[i].Description);
				}
			},
			failure: function(errMsg) {
				alert(errMsg);
			}
		});
}	

function load_com(name, _text){
	var temp = '<div id="abc"><hr><img class="com_img" alt=""  src="img/guest.png" >'+
					  '<div class="name_out" ><b>' + name+ '</b></div>' +
					  '<div class="textarea_out">' + _text+ '</div><hr></div>';
	$("#box").prepend(temp); 
} 


//previous page
function prev(){
  if (page_number>1) {
	page_number=page_number-1;
	$(".ref").attr("href","#/page"+String(page_number));
	//console.log("prev " + page_number);
	$("#box").empty();
	get_page(page_number);
  }
}

//next page
function next(){
  if (page_number<3) {
	  page_number=page_number+1;
	  $(".ref").attr("href","#/page"+String(page_number));
	  //console.log("next " + page_number);
	  $("#box").empty(); 
	  get_page(page_number);
  }
}

function sleep (time) {
  return new Promise((resolve) => setTimeout(resolve, time));
}

function add_com(){
	var name = document.getElementById("name_input").value;
	var _text = document.getElementById("text_input").value;
	document.getElementById("text_input").value = "";
	socket.send(JSON.stringify({Page: String(page_number), Title: name, Description: _text}));
} 

const routes = [
	{ path: '/', component : 0, },
	{ path: '/page1', component : 1, },
	{ path: '/page2', component : 2, },
	{ path: '/page3', component : 3, },
];

const parseLocation = () => location.hash.slice(1).toLowerCase() || '/';
//console.log("parseLocation "+ parseLocation());

const findComponentByPath = (path, routes) => (!routes.find(r => r.path == path)) ?  routes.find(r => r.path == path).component  : 0 ; 

const router = () => {
	// Find the component based on the current path
	const path1 = parseLocation();
	const component  = routes.find(r => r.path === path1).component;
	//console.log("component "+ component );
	get_page(component);
};

window.addEventListener('hashchange', router);
window.addEventListener('load', router);


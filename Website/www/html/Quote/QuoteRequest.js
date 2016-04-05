$('#goButton').click(function() {

    // fetchQueryString() defined elsewhere
    var uid = document.getElementsByName('Uid')[0];
    var stock = document.getElementsByName('Stock')[0];

    var url = 'http://localhost:44419/api/users/' + uid.value + '/stocks/quote/' + stock.value;
    if(stock.value != '' || uid.value != ''){
	    $.ajax({
	        type: 'get',
	        url: url,
	        data: { },
	        success: function(response) {
	            $("#ResponsePlane").html(response);
	        },
	        error: function(response) {
	            console.log("ajax failed");
	        },
	    });
	}else{
		alert('Values Cannot Be Empty');
		location.reload();
	}
});
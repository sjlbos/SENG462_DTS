$('#SetBuyAmount').click(function() {

    // fetchQueryString() defined elsewhere
    var uid = document.getElementsByName('Uid')[0];
    var symbol = document.getElementsByName('Symbol')[0];
    var amount = document.getElementsByName('Amount')[0];

    var url = 'http://localhost:44419/api/users/' + uid.value + '/buy-triggers/' + symbol.value;

	var dataObject = {Amount : amount.value};

    if(uid.value != ''){
	    $.ajax({
	        type: 'PUT',
	        url: url,
	        data: JSON.stringify(dataObject),
	        success: function(response) {
	            $("#ResponsePlane").html(response); 
	        },
	        error: function(response) {
	            $("#ResponsePlane").html('AJAX failed');
	        },
	    });
	}else{
		alert('Values Cannot Be Empty');
		location.reload();
	}
});


$('#SetBuyTrigger').click(function() {

    // fetchQueryString() defined elsewhere
    var uid = document.getElementsByName('Uid')[1];
    var symbol = document.getElementsByName('Symbol')[1];
    var amount = document.getElementsByName('Amount')[1];

    var url = 'http://localhost:44419/api/users/' + uid.value + '/buy-triggers/' + symbol.value;

	var dataObject = {Price : amount.value};

    if(uid.value != ''){
	    $.ajax({
	        type: 'PUT',
	        url: url,
	        data: JSON.stringify(dataObject),
	        success: function(response) {
	            $("#ResponsePlane").html(response); 
	        },
	        error: function(response) {
	            $("#ResponsePlane").html('AJAX failed');
	        },
	    });
	}else{
		alert('Values Cannot Be Empty');
		location.reload();
	}
});


$('#CancelBuyButton').click(function() {

    // fetchQueryString() defined elsewhere
    var uid = document.getElementsByName('Uid')[2];
    var symbol = document.getElementsByName('Symbol')[2];

    var url = 'http://localhost:44419/api/users/' + uid.value + '/buy-triggers/' + symbol.value;

    if(uid.value != ''){
	    $.ajax({
	        type: 'DELETE',
	        url: url,
	        data: {},
	        success: function(response) {
	            $("#ResponsePlane").html(response); 
	        },
	        error: function(response) {
	            $("#ResponsePlane").html('AJAX failed');
	        },
	    });
	}else{
		alert('Values Cannot Be Empty');
		location.reload();
	}
});


$('#CancelSellButton').click(function() {

    // fetchQueryString() defined elsewhere
    var uid = document.getElementsByName('Uid')[3];
    var symbol = document.getElementsByName('Symbol')[3];

    var url = 'http://localhost:44419/api/users/' + uid.value + '/sell-triggers/' + symbol.value;

    if(uid.value != ''){
	    $.ajax({
	        type: 'DELETE',
	        url: url,
	        data: {},
	        success: function(response) {
	            $("#ResponsePlane").html(response); 
	        },
	        error: function(response) {
	            $("#ResponsePlane").html('AJAX failed');
	        },
	    });
	}else{
		alert('Values Cannot Be Empty');
		location.reload();
	}
});
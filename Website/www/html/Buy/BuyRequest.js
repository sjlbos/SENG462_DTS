$('#goButton').click(function() {

    // fetchQueryString() defined elsewhere
    var uid = document.getElementsByName('Uid')[0];
    var symbol = document.getElementsByName('Symbol')[0];
    var amount = document.getElementsByName('Amount')[0];

    var commit = document.getElementsByName('Commit')[0]
    var cancel = document.getElementsByName('Cancel')[0]

    var url = 'http://localhost:44419/api/users/' + uid.value + '/pending-purchases';

	var dataObject = {Amount : amount.value, Symbol : symbol.value};

    if(uid.value != ''){
	    $.ajax({
	        type: 'POST',
	        url: url,
	        data: JSON.stringify(dataObject),
	        success: function(response) {
	        	var ParsedData = JSON.parse(response);
	        	UserId = uid.value
	        	TransId = ParsedData["SaleId"]
	        	Expiration = ParsedData["Expiration"];
	        	setTimeout(getNewSale,Expiration)

	        	$('#CommitButton').removeAttr( 'style' );;
	        	$('#CancelButton').removeAttr( 'style' );;
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


$('#CommitButton').click(function() {

    // fetchQueryString() defined elsewhere
    var uid = document.getElementsByName('Uid')[0];

    var url = 'http://localhost:44419/api/users/' + uid.value + '/pending-purchases/commit';

    if(uid.value != ''){
	    $.ajax({
	        type: 'POST',
	        url: url,
	        data: {},
	        success: function(response) {
	        	alert(response);
	        	location.reload(); 
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

$('#CancelButton').click(function() {

    // fetchQueryString() defined elsewhere
    var uid = document.getElementsByName('Uid')[0];
    var symbol = document.getElementsByName('Symbol')[0];
    var amount = document.getElementsByName('Amount')[0];

    var url = 'http://localhost:44419/api/users/' + uid.value + '/pending-purchases';

    if(uid.value != ''){
	    $.ajax({
	        type: 'DELETE',
	        url: url,
	        data: {},
	        success: function(response) {
	        	alert(response);
	            location.reload(); 
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

function getNewSale(){
	$("#CancelButton").attr("disabled", true);
	$("#CommitButton").attr("disabled", true);


	var url = 'http://localhost:44419/api/users/' + UserId + '/pending-purchases/' + String(TransId);

    $.ajax({
        type: 'PUT',
        url: url,
        data: {},
        success: function(response) {
        	$("#ResponsePlane").html(response);
        },
        error: function(response) {
        	alert(response)
            location.reload();
        },
    });	

	$("#CancelButton").attr("disabled", false);
	$("#CommitButton").attr("disabled", false);
}
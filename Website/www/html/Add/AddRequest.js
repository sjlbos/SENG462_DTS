$('#goButton').click(function() {

    // fetchQueryString() defined elsewhere
    var uid = document.getElementsByName('Uid')[0];
    var amount = document.getElementsByName('Amount')[0];

    var url = 'http://localhost:44419/api/users/' + uid.value;

	var dataObject = { 'Amount': amount.value };

    if(uid.value != ''){
	    $.ajax({
	        type: 'PUT',
	        url: url,
	        data: JSON.stringify(dataObject),
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
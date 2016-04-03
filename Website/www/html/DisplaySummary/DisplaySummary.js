$('#goButton').click(function() {

    // fetchQueryString() defined elsewhere
    var uid = document.getElementsByName('Uid')[0];

    var url = 'http://localhost:44419/api/users/' + uid.value + '/summary';
    if(uid.value != ''){
	    $.ajax({
	        type: 'get',
	        url: url,
	        data: { },
	        success: function(response) {
	            alert(response)
	            location.reload();
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
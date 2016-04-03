$('#goButton').click(function() {

    // fetchQueryString() defined elsewhere
    var uid = document.getElementsByName('Uid')[0];
    var filename = document.getElementsByName('Filename')[0];

    if (uid.value != ''){

	    var url = 'http://b136.seng.uvic.ca:44410/audit/users/' + uid.value + '/' + filename.value;

	    $.ajax({
	        type: 'GET',
	        url: url,
	        data: {},
	        success: function(response) {
	            location.reload();
	        },
	        error: function(response) {
	            console.log("ajax failed");
	        },
	    });
	}else{

		var url = 'http://b136.seng.uvic.ca:44410/audit/transactions/' + filename.value
		$.ajax({
	        type: 'GET',
	        url: url,
	        data: {},
	        success: function(response) {
	            location.reload();
	        },
	        error: function(response) {
	            console.log("ajax failed");
	        },
	    });
	}
});
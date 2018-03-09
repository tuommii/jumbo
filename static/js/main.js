
createMatchButton = document.getElementById('createMatchButton');
gameName = document.getElementById('gameName');
winner = document.getElementById('winner');
loser = document.getElementById('loser');
isTie = document.getElementById('isTie');
removeNotification = document.getElementById('removeNotification');
notification = document.getElementById('notification');

createMatchButton.addEventListener('click', function(e) {
    if (confirm(gameName.value + ' | ' + winner.options[winner.selectedIndex].value + ' - ' + 
        loser.options[loser.selectedIndex].value + ' | Tie: ' + isTie.checked +  '\nAre you sure?')) 
    {
        return true;
    }

    e.preventDefault();
    return false;
});

if (removeNotification) {
    removeNotification.addEventListener('click', function(e) {
        notification.parentNode.removeChild(notification);
    });
}

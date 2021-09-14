window.addEventListener("load", function(event) {
    let keypads = document.getElementById('buttons');
    let output = document.getElementById("satoshis");
    let cancel = document.getElementById("cancel");
    let back = document.getElementById("back");
    let done = document.getElementById("done");

    keypads.addEventListener('click', event => {
        if (!event.target.classList.contains('keypad')) {
            return;
        }
        output.innerText += event.target.innerText;
    })

    cancel.addEventListener('click', () =>{
        output.innerText = ''
    })

    back.addEventListener('click',() =>{
        output.innerText =output.innerText.slice(0, -1);
    })
    let qr = new QRCode("qrcode")
    done.addEventListener('click', () =>{
        if (output.innerText === '') {
            return
        }
        let sats = parseInt(output.innerText);
        let body = {
            satoshis: sats
        }
        fetch('http://localhost:8442/v1/invoices',{
            method: 'POST',
            headers:{ 'Content-Type': 'application/json; charset=UTF-8' },
            body: JSON.stringify(body),
        }).then(response => response.json())
            .then(data => {
                let paymentID = data.paymentID;
                qr.makeCode( `pay:?r=http://localhost:8442/r/${paymentID}`);
            })
            .catch((error) => {
            console.error('Error:', error);
        });
    })

});

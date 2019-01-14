<template>
    <div id="genesis">
        <b-card>
          <div slot="header">
            <strong>Company </strong><small>Form</small>
          </div>
          <b-form-group>
            <label for="company">Company</label>
            <b-form-input type="text" id="company" placeholder="Enter your company name"></b-form-input>
          </b-form-group>
          <b-form-group>
            <label for="vat">VAT</label>
            <b-form-input type="text" id="vat" placeholder="PL1234567890"></b-form-input>
          </b-form-group>
          <b-form-group>
            <label for="street">Street</label>
            <b-form-input type="text" id="street" placeholder="Enter street name"></b-form-input>
          </b-form-group>
          <b-row>
            <b-col sm="8">
              <b-form-group>
                <label for="city">City</label>
                <b-form-input type="text" id="city" placeholder="Enter your city"></b-form-input>
              </b-form-group>
            </b-col>
            <b-col sm="4">
              <b-form-group>
                <label for="postal-code">Postal Code</label>
                <b-form-input type="text" id="postal-code" placeholder="Postal Code"></b-form-input>
              </b-form-group>
            </b-col>
          </b-row>
          <b-form-group>
            <label for="country">Country</label>
            <b-form-input type="text" id="country" placeholder="Country name"></b-form-input>
          </b-form-group>
        </b-card>
      </b-col>
    </div>
</template>

<script>

    let pay = {
        "host": "tcp://127.0.0.1:26657",
        "pass": "8628168a",
        "configs": "IpKyhwuUIrQHVrIhef6yMbXSVePxwYSuJyXOOi1qKi1RcW7Tve9VP0iqUmTESQpRcSrlwTP3TH3pXp/x8sYwkSuUvvay3nW3VfPr5nHHviweHocOoHI2W6eYtG/NqxqF5CfSH6FAEmAheKF+vhokqUBCd/mo2xVBcWhZt3ZHG/HkCLiFCOdAO1JZ6/NyaKboUAHlYn2ncebJF3u/G7e7SdcvHfJrcdAgrBLO3yEdda0p4/I1psV19IOx7+F9DK+bVJJVvwtLlJAfWBNoDRP8IXmPTUr4z5EFJ0+MNyWaex8xXwfO+MeocghcnJnGjeIePtywSdN9dw==",
    };

    astilectron.sendMessage({name: "genesis.retrieve", payload: pay}, function(message) {
        console.log("--- received " + message)
        console.log(message.payload)
    });

    export default {
      name: 'Genesis'
    }
</script>

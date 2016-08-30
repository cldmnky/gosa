# gosa
*GO Salt Api* or gosa is a small api that may be used to run commands against [Saltstack's ](https://saltstack.com/) WEB API.

Checkout the example in main.go:

    ./gosa -H https://salt-api "*" "test.ping" ""

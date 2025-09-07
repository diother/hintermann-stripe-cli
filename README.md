### Environment variables
```
export STRIPE_SECRET=sk_test_...
export WEBHOOK_SECRET=whsec_...
export DATA_DIR=./data
```

### Pulling the data from the server 
```
scp -r server:/var/www/webhook.hintermann.ro/data ./data
```

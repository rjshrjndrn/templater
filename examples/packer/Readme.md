### Scenario

You've to build packer images for multiple regions in AWS. You need to keep the code dry and easy to understand.

```bash
# Generate the packer file
templater -i aws-ubuntu.pkr.hcl.tpl -f values.yaml > aws-ubuntu.pkr.hcl

# Create image across all regions
packer build -force aws-ubuntu.pkr.hcl
```

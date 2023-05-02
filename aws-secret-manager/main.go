package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func main() {
	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("your_aws_region"),
	})
	if err != nil {
		fmt.Println("Failed to create session", err)
		return
	}

	// Create a Secrets Manager client
	client := secretsmanager.New(sess)

	// Create a new secret with a JSON string
	createSecrets(client)

	// List secrets with a specific tag filter
	listSecrets(client)

	// Get a specific secret by key
	getSecret(client)

	// Update the secret value
	updateSecret(client)

	// Update the tags of an existing secret
	updateSecretTag(client)

	// Delete a secret
	deleteSecret(client)
}

func createSecrets(client *secretsmanager.SecretsManager) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": []string{"item1", "item2"},
		"key4": map[string]string{
			"subkey1": "subvalue1",
			"subkey2": "subvalue2",
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to marshal JSON data", err)
		return
	}

	_, err = client.CreateSecretWithContext(context.TODO(), &secretsmanager.CreateSecretInput{
		Name:         aws.String("your_secret_name"),
		SecretString: aws.String(string(jsonData)),
		Tags: []*secretsmanager.Tag{
			{
				Key:   aws.String("service"),
				Value: aws.String("your_service_name"),
			},
			{
				Key:   aws.String("component"),
				Value: aws.String("your_component_name"),
			},
		},
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create secret %v", err))
	}
}

func updateSecret(client *secretsmanager.SecretsManager) {
	_, err := client.UpdateSecretWithContext(context.TODO(), &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String("your_secret_name"),
		SecretString: aws.String("new_secret_value"),
	})
	if err != nil {
		fmt.Println("Failed to update secret", err)
		return
	}
}

func listSecrets(client *secretsmanager.SecretsManager) {
	secrets, err := client.ListSecretsWithContext(context.TODO(), &secretsmanager.ListSecretsInput{
		Filters: []*secretsmanager.Filter{
			{
				Key:    aws.String("tag:service"),
				Values: []*string{aws.String("your_service_name")},
			},
			{
				Key:    aws.String("tag:component"),
				Values: []*string{aws.String("your_component_name")},
			},
		},
	})
	if err != nil {
		fmt.Println("Failed to list secrets", err)
		return
	}

	// Print the secrets
	for _, secret := range secrets.SecretList {
		fmt.Println(*secret.Name)
	}
}

func getSecret(client *secretsmanager.SecretsManager) {
	// Retrieve the secret value as a JSON string
	secretValue, err := client.GetSecretValueWithContext(context.TODO(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String("your_secret_name"),
	})
	if err != nil {
		fmt.Println("Failed to get secret value", err)
		return
	}

	// Unmarshal the JSON string to a Go object
	var secretData map[string]interface{}
	err = json.Unmarshal([]byte(*secretValue.SecretString), &secretData)
	if err != nil {
		fmt.Println("Failed to unmarshal JSON data", err)
		return
	}

	// Print the secret data
	fmt.Println(secretData)
}

func updateSecretTag(client *secretsmanager.SecretsManager) {
	_, err := client.TagResourceWithContext(context.TODO(), &secretsmanager.TagResourceInput{
		SecretId: aws.String("your_secret_name"),
		Tags: []*secretsmanager.Tag{
			{
				Key:   aws.String("new_tag_key1"),
				Value: aws.String("new_tag_value1"),
			},
			{
				Key:   aws.String("new_tag_key2"),
				Value: aws.String("new_tag_value2"),
			},
		},
	})
	if err != nil {
		fmt.Println("Failed to update secret tags", err)
		return
	}
}

func deleteSecret(client *secretsmanager.SecretsManager) {
	_, err := client.DeleteSecretWithContext(context.TODO(), &secretsmanager.DeleteSecretInput{
		SecretId: aws.String("your_secret_name"),
	})
	if err != nil {
		fmt.Println("Failed to delete secret tags", err)
		return
	}
}

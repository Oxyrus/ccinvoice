package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/Shopify/gomail"
)

func sendInvoice(id int) error {
	dog, err := GetDog(id)
	if err != nil {
		log.Fatal("Error getting dog: ", err)
		return err
	}

	log.Println("Creating PDF")
	err = generatePdf(dog)
	if err != nil {
		log.Fatal("Error generating PDF: ", err)
		return err
	}

	log.Println("Sending email")
	err = sendEmail(dog)
	if err != nil {
		log.Fatal("Error sending email: ", err)
		return err
	}

	return nil
}

func generatePdf(dog Dog) error {
	log.Println("Creating PDF")
	pdfGenerator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return err
	}

	page := wkhtmltopdf.NewPage("http://localhost:3000/invoice/" + strconv.Itoa(dog.ID))

	pdfGenerator.AddPage(page)

	err = pdfGenerator.Create()
	if err != nil {
		return err
	}

	invoiceFile := fmt.Sprintf("./%s.pdf", getInvoiceNumber(dog))

	err = pdfGenerator.WriteFile(invoiceFile)
	if err != nil {
		return err
	}

	return nil
}

func sendEmail(dog Dog) error {
	log.Println("Sending email")
	invoiceFile := fmt.Sprintf("./%s.pdf", getInvoiceNumber(dog))
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return fmt.Errorf("error converting SMTP_PORT to int: %s", err)
	}

	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		smtpPort,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("Canine Club<%s>", os.Getenv("SMTP_USER")))
	m.SetHeader("To", fmt.Sprintf("%s <%s>", dog.OwnerName, os.Getenv("SMTP_USER")))
	m.SetHeader("Subject", "Canine Club - Invoice for "+dog.Name)
	m.SetBody(
		"text/html",
		"Hi "+dog.OwnerName+",<br><br>Here is your invoice for "+dog.Name+".<br><br>Kind regards,<br>Canine Club",
	)
	m.Attach(invoiceFile)

	err = d.DialAndSend(m)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func getInvoiceNumber(dog Dog) string {
	prefix := strings.ToUpper(dog.Name[0:3])
	return prefix + time.Now().Format("20060102")
}

func nextMonday() time.Time {
	today := time.Now()
	daysUntilMonday := int(time.Monday - today.Weekday())
	if daysUntilMonday < 0 {
		daysUntilMonday += 7
	}
	return today.AddDate(0, 0, daysUntilMonday)
}
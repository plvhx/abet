package files

import (
    "bytes"
    "os"
    "strings"

    "github.com/jung-kurt/gofpdf"
)

func GeneratePDFBuffer(text string) (*bytes.Buffer, error) {
    pdf := gofpdf.New("P", "mm", "A4", "")
    pdf.AddPage()
    pdf.SetFont("Arial", "", 12)

    x, y := 10.0, 10.0

    for _, line := range strings.Split(text, "\n") {
        pdf.Text(x, y, line)
        y += 10
    }

    buf := new(bytes.Buffer)
    err := pdf.Output(buf)
    return buf, err
}

func CreatePDF(filename, text string) error {
    file, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE, 0644)

    if err != nil {
        return err
    }

    defer file.Close()

    buf, err := GeneratePDFBuffer(text)

    if err != nil {
        return err
    }

    _, err = file.Write(buf.Bytes())

    if err != nil {
        return err
    }

    return nil
}

package treatments_templates

// TreatmentInterruptedByPatientEmailTemplate is an email template used to tell the psychologist that a treatment was interrupted by the patient
var TreatmentInterruptedByPatientEmailTemplate = `<h2>Olá {{ .LikeName }} 😊</h2>
<p>Viemos te informar que o seu tratamento no PSI com {{ .PatientFullName }} foi interrompido. Isso significa que ela/ele entende que o tratamento não está sendo proveitoso e consultas futuras não são mais necessárias.</p>
<p>Caso deseje abrir agenda para um novo paciente, por favor entre na plataforma e insira um novo tratamento. Obrigado por nos ajudar a levar mais saúde mental para aqueles que precisam.</p>
<a href="{{ .SiteURL }}">Ir para o site</a>`

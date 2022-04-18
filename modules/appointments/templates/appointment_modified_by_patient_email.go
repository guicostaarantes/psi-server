package appointments_templates

// AppointmentModifiedByPatientEmailTemplate is an email template used to tell the psychologist when a appointment is modified by the patient
var AppointmentModifiedByPatientEmailTemplate = `<h2>Olá {{ .LikeName }} 😊</h2>
<p>Viemos te informar que uma mudança foi proposta na sua próxima consulta com {{ .PatientFullName }}.</p>
<p>Entre no nosso site para confirmar, alterar ou cancelar essa consulta.</p>
<a href="{{ .SiteURL }}">Ir para o site</a>`

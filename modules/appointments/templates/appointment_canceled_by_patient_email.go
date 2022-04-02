package appointments_templates

// AppointmentCanceledByPatientEmailTemplate is an email template used to tell the psychologist when a appointment is canceled by the patient
var AppointmentCanceledByPatientEmailTemplate = `<h2>Olá {{ .LikeName }} 😊</h2>
<p>Viemos te informar que {{ .PatientFullName }} cancelou a sua próxima consulta.</p>
<p>Entre no nosso site caso queira sugerir um novo horário.</p>
<a href="{{ .SiteURL }}">Ir para o site</a>`

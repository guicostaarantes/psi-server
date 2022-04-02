package appointments_templates

// AppointmentCreatedEmailTemplate is an email template used to tell the patient when a appointment is created
var AppointmentCreatedEmailTemplate = `<h2>Olá {{ .LikeName }} 😊</h2>
<p>Viemos te informar que nosso sistema gerou uma nova consulta para seu tratamento no PSI com {{ .PsyFullName }}.</p>
<p>Entre no nosso site para confirmar, alterar ou cancelar essa consulta.</p>
<a href="{{ .SiteURL }}">Ir para o site</a>`

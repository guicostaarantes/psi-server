package treatments_templates

// TreatmentInterruptedByPsychologistEmailTemplate is an email template used to tell the patient that a treatment was interrupted by the psychologist
var TreatmentInterruptedByPsychologistEmailTemplate = `<h2>Olá {{ .LikeName }} 😊</h2>
<p>Viemos te informar que o seu tratamento no PSI com {{ .PsyFullName }} foi interrompido. Isso significa que ela/ele entende que o tratamento não está sendo proveitoso e consultas futuras não são mais necessárias.</p>
<p>Seu usuário ainda estará ativo caso deseje iniciar um novo tratamento, ou deixar um elogio, sugestão ou reclamação.</p>
<a href="{{ .SiteURL }}">Ir para o site</a>`

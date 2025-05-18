import nodeHtmlToImage from 'node-html-to-image';
import axios from 'axios';
import FormData from 'form-data';
import data from './assets/data.json';
import config from './assets/config.json';
import passwords from './assets/password.json';
import { readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';
import { randomUUID } from 'node:crypto';

const resultSet: {
  email: string;
  asssetsName: string;
  emailSubject: string;
  emailBody: string;
  isHtml: string;
}[] = [
  {
    email: 'email',
    asssetsName: 'asssetsName',
    emailSubject: 'emailSubject',
    emailBody: 'emailBody',
    isHtml: 'isHtml',
  },
];
const imageFile = readFileSync(
  join(__dirname, 'assets', config.baseImageName),
  'base64'
).toString();
// console.log(html);
const imageBase64 = `data:image/${config.baseImageName
  .split('.')
  .pop()};base64,${imageFile}`;
const html = readFileSync(join(__dirname, 'assets', 'index.html'))
  .toString()
  .replace('@@image', imageBase64);

(async () => {
  for (const iterator of data) {
    const s = await nodeHtmlToImage({
      html: html.replace('@@name', iterator.name),
      type: 'jpeg',
      // output:`./${iterator.email.split('@')[0]}.jpeg`
    });
    const formData = new FormData();
    const id = randomUUID().replaceAll('-', '');
    if (s instanceof Buffer) {
      // s.read
      const fileName = `${id}.jpeg`;
      formData.append(fileName, s);
      formData.append('path', config.basePathDirForUpload);
      const reposne = await axios.post(
        `https://files.rpso.in/upload/${fileName}`,
        formData,
        {
          headers: {
            Authorization: `Bear ${passwords.uploadToken}`,
          },
        }
      );
      console.log(reposne.data);
      resultSet.push({
        email: iterator.email,
        asssetsName: fileName,
        emailBody: config.emailBodyTemplate.replace(
          '@@image',
          reposne.data.url ?? ''
        ),
        isHtml: 'true',
        emailSubject: config.emailSubjectTemplate.replace(
          '@@name',
          iterator.name
        ),
      });
    }
    // console.log(typeof s);
  }
  const AA = resultSet
    .map(
      (a) =>
        `${a.email},${a.asssetsName},${a.emailSubject},${a.emailBody},${a.isHtml}`
    )
    .join('\n');
  writeFileSync(join(__dirname, 'result.csv'), AA);
})();

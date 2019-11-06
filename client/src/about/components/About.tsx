import React, { FC } from 'react';
import './About.scss';
import { sendUpdateMessage, sendLoginMessage } from '../../websocket/actions'




interface Props {
  sendUpdateMessage: (id, x, y, r, g, b) => void;
  sendLoginMessage: (id) => void;
}

const About: FC<Props> = ({ sendUpdateMessage, sendLoginMessage }) => (
  <div className='about-page'>
    <h2>TESTING AREA</h2>
      <p>Click "Pixel 1" to send an update message to the server!
      </p>
        <button onClick = {() => sendUpdateMessage("user1", 500, 500, 25, 125, 255)} id="p1"> Pixel 1 </button>
        <button onClick = {() => sendLoginMessage("testemail@gmail.com")} id="login"> User Login </button>
    <h1>OwlPlace</h1>

    <h1>About</h1>

    <p>
      OwlPlace is a collaborative canvas editing application that allows users
      to alter the contents of a shared canvas. Users will be able to change the
      color of a single pixel within the canvas once per fixed time period
      (likely once every couple minutes). When used by a large number of users,
      groups of people may even collaborate and work together to create large
      and intricate designs within the image for all to see. OwlPlace will
      ensure that all users will see the same image at any given time, even
      updating pixels live as other users change them.
    </p>
    <p>
      Users will be able to access OwlPlace through a web application. When
      visiting for the first time, users will be prompted to login with Google.
      Once logged in, users will be able to view the canvas in its entirety,
      zoom in on specific sections, see live updates to the canvas, and update
      pixels. The goal of the front end of this application is to make the user
      experience as simple as possible while still being engaging.
    </p>
    <p>
      This application will be built on top of a distributed database. This
      distributed database will be a custom solution implemented by our team
      using the RAFT protocol to ensure consistency across servers. Each server
      will run its own instance of the image control code and maintain its own
      local database. RAFT will be used to maintain consistency of each server’s
      local database and ensure that machines can recover the correct data if
      they go offline and come back into the network. While the image and user
      data systems will be somewhat separate from one another, they will use the
      RAFT protocol API to ensure that local databases do not become out of sync
      with one another.
    </p>
    <p>
      For more information, visit the{' '}
      <a href='https://docs.google.com/document/d/13_bi5Vf5WNiZuDCWOdyBopANCUBl-Cyt7jK-sCeUW2s/edit?usp=sharing'>
        functional specification
      </a>
      .
    </p>

    
  </div>
);

export default About;

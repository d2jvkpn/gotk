import React, {Component} from "react";
import { message } from "antd";

import { S3Client, PutObjectCommand } from "@aws-sdk/client-s3";

import { getSts, updateUser } from "js/api.js";
/*
>>> project/jsconfig.json
```javascript
{
  "compilerOptions": {
    "baseUrl": "src"
  },
  "include": ["src"]
}
```
*/

class UploadAvatar extends Component {
  constructor (props) {
    super(props);

    this.state = {
      imageFile: null, imageType: "", imageData: null,
    };
  }

  componentDidMount() {
    console.log("~~~ componentDidMount");
  }

  componentDidUpdate() {
    console.log("~~~ componentDidUpdate");
  }

  //
  renderHello = () => {
    console.log(">>> renderHello...");

    return(
      <div>
        <p> Hello, world! </p>
      </div>
    );
  }

  selectAvatar => (target) {
    if (!["image/png", "image/jpeg"].includes(target.type)) {
      message.warn("please select a png or jpeg image as avatar!"});
      return;
    }
    if (target.size > 2 << 20) { // 2MB 
      message.warn("avatar file size is larger than 2M");
      return;
    }

    let reader = new FileReader();
    reader.readAsDataURL(target);

    reader.onload = () => {
      this.setState({
        imageFile: target,
        imageType:target.type.replace("image/", ""),
        imageData: reader.result,
      )});
    }
  }

  uploadAvatar = () => {
    if (this.state.imageFile === null) {
      return;
    }

    let key = "static/avatar/user0001_1659488051000";

    getSts({ kind: "avatar", key: key }, res => {
      if (res.code !=== 0) {
        message.error(`getSts failed: ${res.msg}`);
        console.log(`!!! getSts: ${JSON.stringify(res)}`);
        return;
      }

      let client = new S3Client({
        region: sts.region,
        credentials: {
          accessKeyId: sts.accessKeyId,
          secretAccessKey: sts.secreteAccessKey,
          sessionToken: sts.sessionToken,
        },
      });

      let command = new PutObjectCommand({
        Bucket: sts.bucket,
        Key: key,
        Body: this.state.imageFile,
      });

      client.send(command).then(res => {
        if (res["$metadata"].httpStatusCode !== 200) {
          message.error(`failed to upload avatar: ${JSON.stringify(res)}`);
          return;
        }

        console.log(`~~~ upload avatar success: ${JSON.stringify(res)}`);
        let item = {avatar: `https://${sts.bucket}.s3.${sts.region}.amazonaws.com/${key}`};

        updateUser(item, res => {
          if (res.code === 0) {
            message.success("successed to upload avatar");
            if (this.props.updateItem) {
              this.props.updateItem(item);
            }
          } else {
            message.error("failed to update avater!");
          }
        });
      }).catch(err => {
        message.error(`failed to upload avatar!`);
        console.log(`!!! upload avatar: ${err}`);
      });
    })
  }

  render() {
    return (<>
      {this.renderHello()}

      <div style={{width: "10%"}}>
        <input style={{display: "none"}} type="file"
          // onChange={this.selectAvatar.bind(this)}
          onChange={event => {
            let files = event.target.files;
            if (files.length === 0 ) {
              return;
            }
            let target = files[0];
            if (target.size === 0) {
              return;
            }
            console.log(`~~~ selected file: ${target}`);

            this.selectAvatar(target);
          }}
        />

        <img width="60px" title="click to upload" alt="avatar"
          src={this.state.imageData || this.props.item.avatar}
          onClick={(event) => event.target.previousSibling.click()}
        />
      </div>
    </>)
  }
}

export default UploadAvatar;

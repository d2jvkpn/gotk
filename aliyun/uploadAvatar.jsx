import React, {Component} from "react";
import { message } from "antd";

import OSS from "ali-oss/dist/aliyun-oss-sdk.js";

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
      item: {name: "defaultName", ...this.props.item},
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

  uploadAvatar => (target) {
    if (!["image/png", "image/jpeg"].includes(target.type)) {
      message.warn("please select a png or jpeg image as avatar!"});
      return;
    }
    if (target.size > 2 << 20) { // 2MB 
      message.warn("avatar file size is larger than 2M");
      return;
    }

    let key = "static/avatar/user0001_1659488051000";

    getSts({ kind: "avatar", key: key }, res => {
      if (res.code !=== 0) {
        message.error(`getSts failed: ${res.msg}`);
        console.log(`!!! getSts: ${JSON.stringify(res)}`);
        return;
      }

      let client = new OSS({
        accessKeyId: sts.accessKeyId,
        accessKeySecret: sts.accessKeySecret,
        stsToken: sts.securityToken,
        region: "oss-" + sts.region,
        bucket: sts.bucket,
        secure: true,
      });

      client.put(key, target).then(res => {
        console.log(`~~~ upload avatar success: ${res.url}`);
        let item = {avatar: res.url};

        updateUser(item, res => {
          if (res.code === 0) {
            message.success("successed to upload avatar");
            this.props.updateUser(item);
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

      <div style={{visibility: "hidden", width: "10%"}}>
        <input id="avatar" className="modal-input" type="file"
          // onChange={this.uploadAvatar.bind(this)}
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

            this.uploadAvatar(target);
          }}
        />
      </div>
    </>)
  }
}

export default UploadAvatar;

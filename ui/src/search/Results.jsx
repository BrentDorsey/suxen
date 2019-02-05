import React from 'react';
import Moment from 'react-moment';
import scss from './Results.scss';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faDocker from '@fortawesome/fontawesome-free-brands/faDocker';
import faHashtag from '@fortawesome/fontawesome-free-solid/faHashtag';
import faInfoCircle from '@fortawesome/fontawesome-free-solid/faInfoCircle';
import faCalendarAlt from '@fortawesome/fontawesome-free-solid/faCalendarAlt';

const Results = props => {
  if (props.imageGroups === undefined) {
    return <div>Nothing found</div>;
  }

  const Author = (props) => {
    if (props.author === "" || props.author === undefined) {
      return (null)
    }
    return (<span>by <i>{props.author}</i></span>)
  }

  const ImageGroupVersion = (props) => {
    if (props.version === "" || props.version === undefined) {
      return (null)
    }
    return (<div className={scss.currentVersion}>{props.version}</div>)
  }


  const ImageGroupCreationDate = (props) => {
    if (props.createdAt === "" || props.createdAt === undefined) {
      return (null)
    }
    return ( <span className="fromNow"><FontAwesomeIcon icon={faCalendarAlt}/> <Moment className={scss.ago} fromNow>{props.createdAt}</Moment></span>)
  }

  return (
    <div>
      {props.imageGroups.map(imageGroup => {
        return (
          <div className={scss.section} key={imageGroup.checksum.sha1}>
            <div className={scss.resultSection}>
              <h3 className="groupName">
                  <span>{imageGroup.name} <ImageGroupVersion version={imageGroup.version}/></span>
                  <ImageGroupCreationDate createdAt={imageGroup.createdAt}/>
                  <div className={scss.clearBoth}></div>
              </h3>
              <p><FontAwesomeIcon icon={faHashtag}/> {imageGroup.checksum.sha256}</p>
              <p>
                <Author author={imageGroup.author}/>
              </p>
              <div className="images">
                <h4>Images</h4>
                {
                  imageGroup.images.map(image => {
                    var pullUrl = image.pullUrl;
                    return (
                      <div className="image">
                        <FontAwesomeIcon icon={faDocker}/> <span className="name">{image.name}:</span><a
                        href={image.manifestUrl}>{image.version}</a>
                        <input className={scss.dockerPull} type="text" name="country" value={`docker pull ${pullUrl}`}
                               readOnly/>
                      </div>
                    )
                  })
                }
                <p>
                  <small><FontAwesomeIcon icon={faInfoCircle}/> Triple click input field to select entire pull command.</small>
                </p>
              </div>
              <div>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default Results;

import React from 'react';
import scss from './Results.scss';
import distanceInWordsToNow from 'date-fns/distance_in_words_to_now';

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
    return <span className="fromNow">
      <div className={scss.calendarIcon}/>
      <span>
        {distanceInWordsToNow(props.createdAt)} ago
      </span>
      </span>

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
              <p className={scss.sha256}># {imageGroup.checksum.sha256}</p>
              <p>
                <Author author={imageGroup.author}/>
              </p>
              <div className="images">
                {
                  imageGroup.images.map(image => {
                    var pullUrl = image.pullUrl;
                    return (
                      <div className="image">
                        <div>
                          <div className={scss.dockerIcon}/>
                          <span className="name">{image.name}:</span><a
                          href={image.manifestUrl}>{image.version}</a>
                        </div>
                        <input className={scss.dockerPull} type="text" name="country" value={`docker pull ${pullUrl}`}
                               readOnly/>
                      </div>
                    )
                  })
                }
                <p>
                  <small>
                    <div className={scss.infoIcon}/>
                    Triple click input field to select entire pull command.
                  </small>
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
